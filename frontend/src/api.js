// API Configuration
const API_URL = 'http://localhost:8080/api';

export async function registerUser(email, password) {
    try {
        console.log('API: Registering user:', email);
        const response = await fetch(`${API_URL}/auth/register`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ email, password })
        });

        console.log('API: Register response status:', response.status);

        const data = await response.json();
        
        if (!response.ok) {
            throw new Error(data.error || 'Registration failed');
        }

        return data;
    } catch (error) {
        console.error('API: Register error:', error);
        throw error;
    }
}

export async function loginUser(email, password) {
    try {
        console.log('API: Logging in user:', email);
        const response = await fetch(`${API_URL}/auth/login`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ email, password })
        });

        console.log('API: Login response status:', response.status);

        const data = await response.json();
        
        if (!response.ok) {
            throw new Error(data.error || 'Login failed');
        }

        console.log('API: Login success, token received');
        return data;
    } catch (error) {
        console.error('API: Login error:', error);
        throw error;
    }
}

function fetchWithTimeout(resource, options = {}) {
    const { timeout = 15000 } = options
    const controller = new AbortController()
    const id = setTimeout(() => controller.abort(), timeout)
    options.signal = controller.signal

    return fetch(resource, options).finally(() => clearTimeout(id))
}

async function doFetch(url, opts = {}) {
    console.log('[api] FETCH ->', url, opts && opts.method ? opts.method : 'GET')
    try {
        const res = await fetchWithTimeout(url, opts)
        console.log('[api] RESPONSE <-', url, res && res.status)
        if (!res.ok) {
            const text = await res.text().catch(() => '')
            console.error('[api] RESPONSE ERROR', url, res.status, text)
            // If unauthorized, clear client auth and redirect to login
            if (res.status === 401) {
                try {
                    console.warn('[api] Unauthorized - clearing local session')
                    localStorage.removeItem('token')
                    localStorage.removeItem('role')
                    localStorage.removeItem('user')
                    localStorage.removeItem('archive')
                } catch (e) {
                    console.error('[api] error clearing storage', e)
                }
                // redirect user to login page
                if (typeof window !== 'undefined') {
                    window.location.href = './index.html'
                }
            }
            throw new Error(text || `HTTP ${res.status}`)
        }
        const json = await res.json().catch(() => null)
        console.log('[api] PARSED JSON <-', url, json ? (Array.isArray(json) ? `array(${json.length})` : 'object') : 'null')
        return json
    } catch (err) {
        if (err.name === 'AbortError') console.error('[api] FETCH TIMEOUT', url)
        else console.error('[api] FETCH ERROR', url, err)
        throw err
    }
}

export async function getModels(archiveId) {
    const token = localStorage.getItem('token');
    const headers = {};
    if (token) headers['Authorization'] = `Bearer ${token}`;
    let url = `${API_URL}/models`;
    if (archiveId) url += `?archive_id=${archiveId}`;
    const data = await doFetch(url, { headers });
    return (data && data.data) || [];
}

export async function uploadModel(file, name, description) {
    const token = localStorage.getItem('token');
    const formData = new FormData();
    formData.append('file', file);
    formData.append('name', name);
    formData.append('description', description);
    const archiveId = document.getElementById('archiveSelect') ? document.getElementById('archiveSelect').value : '';
    if (archiveId) formData.append('archive_id', archiveId);

    const response = await fetch(`${API_URL}/models/upload`, {
        method: 'POST',
        headers: {
            'Authorization': `Bearer ${token}`
        },
        body: formData
    });

    if (!response.ok) {
        const error = await response.json();
        throw new Error(error.error || 'Upload failed');
    }

    return await response.json();
}

export async function createArchive(name) {
    const token = localStorage.getItem('token');
    const response = await fetch(`${API_URL}/archives`, {
        method: 'POST',
        headers: {
            'Authorization': `Bearer ${token}`
        },
        body: new URLSearchParams({ name: name || '' })
    });
    if (!response.ok) {
        const err = await response.json();
        throw new Error(err.error || 'Failed to create archive');
    }
    return await response.json();
}

export async function listArchives() {
    const token = localStorage.getItem('token');
    const headers = {};
    if (token) headers['Authorization'] = `Bearer ${token}`;
    return await doFetch(`${API_URL}/archives`, { headers });
}

export async function deleteArchive(id) {
    const token = localStorage.getItem('token');
    const response = await fetch(`${API_URL}/archives`, {
        method: 'DELETE',
        headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ id })
    });
    if (!response.ok) {
        const err = await response.json();
        throw new Error(err.error || 'Failed to delete archive');
    }
    return await response.json();
}

export async function archiveLogin(tokenStr) {
    const response = await fetch(`${API_URL}/archives/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ token: tokenStr })
    });
    if (!response.ok) {
        const err = await response.json();
        throw new Error(err.error || 'Login failed');
    }
    return await response.json();
}

export async function getUserProfile() {
    const token = localStorage.getItem('token');
    const response = await fetch(`${API_URL}/user/profile`, {
        headers: {
            'Authorization': `Bearer ${token}`
        }
    });

    if (!response.ok) {
        throw new Error('Failed to fetch user profile');
    }

    return await response.json();
}
