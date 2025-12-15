import { loginUser, registerUser, archiveLogin } from './api.js';

window.toggleForm = function() {
    const loginForm = document.getElementById('loginForm');
    const registerForm = document.getElementById('registerForm');
    loginForm.style.display = loginForm.style.display === 'none' ? 'block' : 'none';
    registerForm.style.display = registerForm.style.display === 'none' ? 'block' : 'none';
};

function showMessage(msg, type) {
    const messageDiv = document.getElementById('message');
    messageDiv.textContent = msg;
    messageDiv.className = `message ${type}`;
    messageDiv.style.display = 'block';
    setTimeout(() => {
        messageDiv.style.display = 'none';
    }, 5000);
}

// Login form
const loginForm = document.getElementById('form-login');
if (loginForm) {
    loginForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        console.log('Login form submitted');
        
        const email = document.getElementById('email').value;
        const password = document.getElementById('password').value;
        
        console.log('Attempting login:', email);
        
        try {
            const result = await loginUser(email, password);
            console.log('Login successful:', result);
            
            localStorage.setItem('token', result.token);
            localStorage.setItem('user', JSON.stringify(result.user));
            localStorage.setItem('role', result.role);
            
            showMessage('Login berhasil!', 'success');
            
            // Redirect based on role
            setTimeout(() => {
                if (result.role === 'admin') {
                    window.location.href = './admin.html';
                } else {
                    window.location.href = './viewer.html';
                }
            }, 1000);
        } catch (error) {
            console.error('Login error:', error);
            showMessage('Login gagal: ' + error.message, 'error');
        }
    });
}

// Register form
const registerForm = document.getElementById('form-register');
if (registerForm) {
    registerForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        console.log('Register form submitted');
        
        const email = document.getElementById('reg-email').value;
        const password = document.getElementById('reg-password').value;
        
        console.log('Attempting register:', email);
        
        try {
            const result = await registerUser(email, password);
            console.log('Register successful:', result);
            
            showMessage('Registrasi berhasil! Silakan login.', 'success');
            setTimeout(() => {
                window.toggleForm();
            }, 2000);
        } catch (error) {
            console.error('Register error:', error);
            showMessage('Registrasi gagal: ' + error.message, 'error');
        }
    });
}

// Check if already logged in
window.addEventListener('DOMContentLoaded', () => {
    console.log('Page loaded');
    if (localStorage.getItem('token')) {
        console.log('Token found, redirecting...');
        const role = localStorage.getItem('role');
        if (role === 'admin') {
            window.location.href = './admin.html';
        } else {
            window.location.href = './viewer.html';
        }
    }
    // Slap toggle behavior
    const slapToggle = document.getElementById('slapToggle');
    const standardArea = document.getElementById('standardArea');
    const tokenArea = document.getElementById('tokenArea');
    if (slapToggle) {
        slapToggle.addEventListener('click', (e) => {
            const btn = e.target.closest('.slap-option');
            if (!btn) return;
            // remove active from all
            slapToggle.querySelectorAll('.slap-option').forEach(x => x.classList.remove('active'));
            btn.classList.add('active');
            const mode = btn.dataset.mode;
            if (mode === 'guest') {
                standardArea.style.display = 'none';
                tokenArea.style.display = 'block';
            } else {
                standardArea.style.display = 'block';
                tokenArea.style.display = 'none';
            }
        });
    }
    const tokenBtn = document.getElementById('tokenLoginBtn');
    if (tokenBtn) {
        tokenBtn.addEventListener('click', async () => {
            const val = document.getElementById('archiveTokenInput').value.trim();
            if (!val) { showMessage('Masukkan token', 'error'); return; }
            try {
                const res = await archiveLogin(val);
                // backend returns { token, archive }
                localStorage.setItem('token', res.token);
                localStorage.setItem('role', 'archive_user');
                localStorage.setItem('archive', JSON.stringify(res.archive));
                showMessage('Login dengan token berhasil', 'success');
                setTimeout(() => { window.location.href = './viewer.html'; }, 800);
            } catch (err) {
                console.error('Token login failed', err);
                showMessage('Login token gagal: ' + err.message, 'error');
            }
        });
    }
});

