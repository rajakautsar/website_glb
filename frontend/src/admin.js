import { getModels, listArchives, createArchive, deleteArchive } from './api.js';

// Check authorization
function checkAuth() {
    const token = localStorage.getItem('token');
    const role = localStorage.getItem('role');
    
    if (!token || role !== 'admin') {
        window.location.href = './index.html';
        return;
    }

    const user = JSON.parse(localStorage.getItem('user'));
    document.getElementById('userEmail').textContent = user.email;
}

// Load models
let _isLoadingModels = false;
async function loadModels() {
    if (_isLoadingModels) return;
    _isLoadingModels = true;
    try {
        const models = await getModels();
        displayModels(models);
    } catch (error) {
        console.error('Error loading models:', error);
    } finally {
        _isLoadingModels = false;
    }
}

function displayModels(models) {
    const container = document.getElementById('modelsList');
    if (models.length === 0) {
        container.innerHTML = '<p>Belum ada model</p>';
        return;
    }

    container.innerHTML = models.map(model => `
        <div class="model-card">
            <h3>${model.name}</h3>
            <p>${model.description || 'No description'}</p>
            <p class="model-info">Upload: ${model.uploaded_by}</p>
            <p class="model-info">Size: ${(model.file_size / 1024).toFixed(2)} KB</p>
            <button onclick="viewModel(${model.id})" class="btn btn-small">View</button>
            <button onclick="deleteModel(${model.id})" class="btn btn-danger btn-small">Delete</button>
        </div>
    `).join('');
}

window._deleteModelFn = async function(id) {
    if (!confirm('Delete model ini?')) return;
    console.log('Deleting model (fn) ->', id);
    try {
        const token = localStorage.getItem('token');
        const controller = new AbortController();
        const timeoutId = setTimeout(() => controller.abort(), 15000);
        const response = await fetch(`http://localhost:8080/api/models`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ id: id }),
            signal: controller.signal
        });
        clearTimeout(timeoutId);
        console.log('Delete response status:', response.status);
        if (response.ok) {
            showMessage('Model deleted', 'success');
            loadModels();
            // refresh archives and open archive view if needed
            loadArchives();
            if (currentOpenArchiveId) await loadArchiveFiles(currentOpenArchiveId);
        } else {
            const err = await response.text().catch(() => '');
            console.error('Delete failed body:', err);
            showMessage('Error deleting model', 'error');
        }
    } catch (error) {
        if (error.name === 'AbortError') console.error('Delete request aborted (timeout)')
        else console.error('Error deleting model (fn):', error);
        showMessage('Error deleting model', 'error');
    }
};

// Handle upload
document.getElementById('uploadForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const name = document.getElementById('modelName').value;
    const description = document.getElementById('modelDescription').value;
    const file = document.getElementById('modelFile').files[0];

    try {
        const token = localStorage.getItem('token');
        const formData = new FormData();
        formData.append('file', file);
        formData.append('name', name);
        formData.append('description', description);

        const archiveId = document.getElementById('archiveSelect').value || '';
        const response = await fetch('http://localhost:8080/api/models/upload', {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${token}`
            },
            body: (function(){ formData.append('archive_id', archiveId); return formData })()
        });

        if (response.ok) {
            showMessage('Model uploaded successfully!', 'success');
            document.getElementById('uploadForm').reset();
            loadModels();
        } else {
            const error = await response.json();
            showMessage('Upload failed: ' + error.error, 'error');
        }
    } catch (error) {
        showMessage('Upload error: ' + error.message, 'error');
    }
});

// Archive management
async function loadArchives() {
    try {
        const res = await listArchives();
        displayArchives(res.data || []);
    } catch (err) {
        console.error('Error loading archives', err);
        document.getElementById('archivesList').innerHTML = '<p>Failed to load archives</p>';
    }
}

function displayArchives(arr) {
    const container = document.getElementById('archivesList');
    const select = document.getElementById('archiveSelect');
    select.innerHTML = '<option value="">(None)</option>';
    if (!arr || arr.length === 0) {
        container.innerHTML = '<p>Belum ada arsip</p>';
        return;
    }

    // build DOM nodes so we can attach open/show files behavior
    container.innerHTML = '';
    arr.forEach(a => {
        // archive card
        const card = document.createElement('div');
        card.className = 'model-card';
        card.innerHTML = `
            <div class="archive-header">
                <div>
                    <h3 style="margin:0">${a.name}</h3>
                    <p style="margin:4px 0; color:var(--text-secondary)">TOKEN: <code>${a.token}</code></p>
                    <p style="margin:4px 0; color:var(--text-secondary)">Jumlah model: ${a.count}</p>
                </div>
            </div>
            <div class="archive-files" id="archive-files-${a.id}" style="display:none; margin-top:12px;"></div>
        `;

        // make header clickable to toggle files
        const header = card.querySelector('.archive-header');
        header.style.cursor = 'pointer';
        header.addEventListener('click', async () => {
            const filesContainer = document.getElementById(`archive-files-${a.id}`);
            if (filesContainer.style.display === 'none') {
                filesContainer.style.display = 'block';
                await loadArchiveFiles(a.id);
            } else {
                filesContainer.style.display = 'none';
            }
        });

        container.appendChild(card);

        // populate select
        const opt = document.createElement('option');
        opt.value = a.id;
        opt.textContent = a.name;
        select.appendChild(opt);
    });
}

let currentOpenArchiveId = null;
async function loadArchiveFiles(archiveId) {
    const filesContainer = document.getElementById(`archive-files-${archiveId}`);
    if (!filesContainer) return;
    filesContainer.innerHTML = '<p>Loading files...</p>';
    try {
        const models = await getModels(archiveId);
        if (!models || models.length === 0) {
            filesContainer.innerHTML = '<p>Tidak ada model di arsip ini</p>';
            return;
        }
        filesContainer.innerHTML = '';
        models.forEach(m => {
            const row = document.createElement('div');
            row.className = 'model-row';
            row.style.display = 'flex';
            row.style.justifyContent = 'space-between';
            row.style.alignItems = 'center';
            row.style.padding = '6px 0';
            row.style.borderBottom = '1px solid var(--border-color)';

            const info = document.createElement('div');
            const title = document.createElement('strong');
            title.textContent = m.name;
            const meta = document.createElement('div');
            meta.style.fontSize = '0.85rem';
            meta.style.color = 'var(--text-secondary)';
            meta.textContent = `${(m.file_size/1024).toFixed(2)} KB`;
            info.appendChild(title);
            info.appendChild(meta);

            const controls = document.createElement('div');
            controls.style.display = 'flex';
            controls.style.gap = '8px';
            controls.style.alignItems = 'center';

            const viewBtn = document.createElement('button');
            viewBtn.className = 'btn btn-small';
            viewBtn.textContent = 'View';
            viewBtn.addEventListener('click', (e) => { e.stopPropagation && e.stopPropagation(); viewModel(m.id); });

            const delBtn = document.createElement('button');
            delBtn.className = 'btn btn-danger btn-small';
            delBtn.textContent = 'Delete';
            delBtn.addEventListener('click', async (e) => {
                // stop propagation to avoid collapsing the panel
                e.stopPropagation();
                if (!confirm('Delete model ini?')) return;
                try {
                    // disable buttons to prevent double click
                    delBtn.disabled = true;
                    viewBtn.disabled = true;
                    console.log('Deleting model', m.id);
                    await window._deleteModelFn(m.id);
                } catch (err) {
                    console.error('Error deleting model (client):', err);
                } finally {
                    delBtn.disabled = false;
                    viewBtn.disabled = false;
                }
            });

            controls.appendChild(viewBtn);
            controls.appendChild(delBtn);

            row.appendChild(info);
            row.appendChild(controls);
            filesContainer.appendChild(row);
        });

        // footer actions inside the archive panel
        const footer = document.createElement('div');
        footer.style.display = 'flex';
        footer.style.justifyContent = 'flex-end';
        footer.style.gap = '8px';
        footer.style.marginTop = '12px';

        const closeBtn = document.createElement('button');
        closeBtn.className = 'btn';
        closeBtn.textContent = 'Tutup';
        closeBtn.addEventListener('click', () => { filesContainer.style.display = 'none'; });

        const delArchiveBtn = document.createElement('button');
        delArchiveBtn.className = 'btn btn-danger';
        delArchiveBtn.textContent = 'Hapus Folder Arsip';
        delArchiveBtn.addEventListener('click', async (e) => {
            e.stopPropagation && e.stopPropagation();
            if (!confirm('Hapus folder arsip beserta isinya?')) return;
            try {
                delArchiveBtn.disabled = true;
                console.log('Deleting archive', archiveId);
                await window._deleteArchive(archiveId);
            } catch (err) {
                console.error('Error deleting archive (client):', err);
            } finally {
                delArchiveBtn.disabled = false;
            }
        });

        footer.appendChild(closeBtn);
        footer.appendChild(delArchiveBtn);
        filesContainer.appendChild(footer);

        currentOpenArchiveId = archiveId;
    } catch (err) {
        console.error('Failed to load archive files', err);
        filesContainer.innerHTML = '<p>Gagal memuat file</p>';
    }
}

document.getElementById('createArchiveBtn').addEventListener('click', async () => {
    if (!confirm('Buat folder arsip baru?')) return;
    try {
        await createArchive('');
        showMessage('Archive created', 'success');
        loadArchives();
        loadModels();
    } catch (err) {
        showMessage('Failed to create archive', 'error');
    }
});

window._deleteArchive = async function(id) {
    if (!confirm('Delete archive ini? This will remove all models inside.')) return;
    try {
        await deleteArchive(id);
        showMessage('Archive deleted', 'success');
        loadArchives();
        loadModels();
    } catch (err) {
        showMessage('Failed to delete archive', 'error');
    }
}

function showMessage(msg, type) {
    const messageDiv = document.getElementById('uploadMessage');
    messageDiv.textContent = msg;
    messageDiv.className = `message ${type}`;
    messageDiv.style.display = 'block';
    setTimeout(() => {
        messageDiv.style.display = 'none';
    }, 3000);
}

// Initialize
checkAuth();
loadModels();
loadArchives();
setInterval(loadModels, 5000);
