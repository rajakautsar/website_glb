import * as THREE from 'three';
import { GLTFLoader } from 'three/examples/jsm/loaders/GLTFLoader.js';
import { OrbitControls } from 'three/examples/jsm/controls/OrbitControls.js';
import { getModels } from './api.js';

let scene, camera, renderer, controls, currentModel;
let models = [];
let raycaster, pointer, selectionBox;
let isAnimating = false;
let prevCameraPos = null;
let prevTarget = null;

// Check auth
function checkAuth() {
    const token = localStorage.getItem('token');
    if (!token) {
        window.location.href = './index.html';
        return;
    }
    const role = localStorage.getItem('role');
    if (role === 'admin' || role === 'user') {
        const userStr = localStorage.getItem('user');
        if (!userStr) { window.location.href = './index.html'; return; }
        const user = JSON.parse(userStr);
        const el = document.getElementById('userEmail');
        if (el) el.textContent = user.email;
    } else if (role === 'archive_user') {
        const archStr = localStorage.getItem('archive');
        if (!archStr) { window.location.href = './index.html'; return; }
        const arch = JSON.parse(archStr);
        const el = document.getElementById('userEmail');
        if (el) el.textContent = arch.name || arch.token || 'Archive User';
    } else {
        // unknown role -> redirect
        window.location.href = './index.html';
        return;
    }
}

// Initialize Three.js
function initThreeJS() {
    const container = document.getElementById('canvas-container');
    
    console.log('Container dimensions:', {
        clientWidth: container.clientWidth,
        clientHeight: container.clientHeight,
        offsetWidth: container.offsetWidth,
        offsetHeight: container.offsetHeight
    });
    
    scene = new THREE.Scene();
    // white background for the viewer canvas
    scene.background = new THREE.Color(0xffffff);
    
    const width = container.clientWidth || window.innerWidth * 0.8;
    const height = container.clientHeight || window.innerHeight * 0.8;
    
    // use small near plane to avoid clipping when zooming in
    camera = new THREE.PerspectiveCamera(
        75,
        width / height,
        0.01,
        100000
    );
    camera.position.set(0, 0, 10);

    renderer = new THREE.WebGLRenderer({ 
        antialias: true,
        alpha: false,
        preserveDrawingBuffer: true
    });
    renderer.setSize(width, height);
    renderer.setPixelRatio(window.devicePixelRatio);
    renderer.shadowMap.enabled = true;
    // ensure renderer clears to white
    if (renderer.setClearColor) renderer.setClearColor(0xffffff, 1);
    renderer.outputColorSpace = THREE.SRGBColorSpace;
    
    // Make sure container is cleared and canvas is added
    container.innerHTML = '';
    container.appendChild(renderer.domElement);
    
    console.log('Renderer created with size:', width, 'x', height);

    // Lighting - more intense
    const ambientLight = new THREE.AmbientLight(0xffffff, 0.9);
    scene.add(ambientLight);

    const directionalLight = new THREE.DirectionalLight(0xffffff, 1.0);
    directionalLight.position.set(10, 20, 10);
    directionalLight.castShadow = true;
    directionalLight.shadow.mapSize.width = 2048;
    directionalLight.shadow.mapSize.height = 2048;
    scene.add(directionalLight);
    
    console.log('Lights added');

    // Controls
    controls = new OrbitControls(camera, renderer.domElement);
    controls.enableDamping = true;
    controls.dampingFactor = 0.05;
    controls.autoRotate = false;
    // prevent panning from moving camera off-center
    controls.screenSpacePanning = false;
    // sensible defaults for min/max distances (will be adjusted per-model)
    controls.minDistance = 0.1;
    controls.maxDistance = 1000;

    // Animation loop
    function animate() {
        requestAnimationFrame(animate);
        controls.update();
        renderer.render(scene, camera);
    }
    animate();
    
    // setup raycaster for picking (double-click to focus)
    raycaster = new THREE.Raycaster();
    pointer = new THREE.Vector2();
    renderer.domElement.addEventListener('dblclick', onDoubleClick);
    
    console.log('Animation loop started');

    // Handle resize
    window.addEventListener('resize', () => {
        const w = container.clientWidth || window.innerWidth * 0.8;
        const h = container.clientHeight || window.innerHeight * 0.8;
        camera.aspect = w / h;
        camera.updateProjectionMatrix();
        renderer.setSize(w, h);
    });
}

// ----- Picking & focus helpers -----
function onDoubleClick(event) {
    if (!currentModel || isAnimating) return;
    const rect = renderer.domElement.getBoundingClientRect();
    pointer.x = ((event.clientX - rect.left) / rect.width) * 2 - 1;
    pointer.y = -((event.clientY - rect.top) / rect.height) * 2 + 1;

    raycaster.setFromCamera(pointer, camera);
    const intersects = raycaster.intersectObject(currentModel, true);
    if (intersects.length > 0) {
        const picked = intersects[0].object;
        focusOnObject(picked);
    } else {
        // clicked empty space: cancel selection and return to previous view if available
        if (selectionBox) {
            try { scene.remove(selectionBox); } catch (e) {}
            selectionBox = null;
        }
        if (prevCameraPos && prevTarget) {
            // animate back to previous camera/target
            animateCameraTo(prevCameraPos.clone(), prevTarget.clone(), 500);
            // clear stored previous once animating back
            prevCameraPos = null;
            prevTarget = null;
        }
    }
}

function focusOnObject(object) {
    if (!object) return;
    // store previous camera/target to allow cancel/back
    try {
        prevCameraPos = camera.position.clone();
        prevTarget = controls.target.clone();
    } catch (e) {
        prevCameraPos = null;
        prevTarget = null;
    }
    // build bounding box for the picked object
    const box = new THREE.Box3().setFromObject(object);
    const center = box.getCenter(new THREE.Vector3());
    const size = box.getSize(new THREE.Vector3());
    const maxDim = Math.max(size.x, size.y, size.z);

    // compute a suitable camera distance
    const fov = camera.fov * (Math.PI / 180);
    let cameraZ = Math.abs(maxDim / 2 / Math.tan(fov / 2));
    cameraZ = Math.max(cameraZ * 1.2, maxDim * 1.2, 0.2);

    // direction from center to current camera
    const dir = camera.position.clone().sub(controls.target || new THREE.Vector3()).normalize();
    const targetPos = center.clone().add(dir.multiplyScalar(cameraZ));

    // adjust camera near/far and controls ranges
    camera.near = Math.min(0.01, Math.max(0.001, maxDim / 1000));
    camera.far = Math.max(1000, cameraZ * 20);
    camera.updateProjectionMatrix();

    const targetLookAt = center.clone();

    // show selection box helper
    if (selectionBox) {
        scene.remove(selectionBox);
        selectionBox.geometry && selectionBox.geometry.dispose();
    }
    try {
        selectionBox = new THREE.BoxHelper(object, 0xff0000);
        scene.add(selectionBox);
    } catch (e) {
        console.warn('Could not add BoxHelper:', e);
    }

    // animate camera + controls.target
    animateCameraTo(targetPos, targetLookAt, 600);
}

function animateCameraTo(targetPos, targetLookAt, duration = 600) {
    isAnimating = true;
    const startPos = camera.position.clone();
    const startTarget = controls.target.clone();
    const startTime = performance.now();

    function tick(now) {
        const t = Math.min(1, (now - startTime) / duration);
        const ease = t < 0.5 ? 2 * t * t : -1 + (4 - 2 * t) * t; // easeInOut

        camera.position.lerpVectors(startPos, targetPos, ease);
        const newTarget = startTarget.clone().lerp(targetLookAt, ease);
        controls.target.copy(newTarget);
        controls.update();

        if (t < 1) {
            requestAnimationFrame(tick);
        } else {
            // finalize
            camera.position.copy(targetPos);
            controls.target.copy(targetLookAt);
            controls.update();
            isAnimating = false;
        }
    }

    requestAnimationFrame(tick);
}

// Load models list
async function loadModelsList() {
    try {
        console.log('Loading models list...');
        models = await getModels();
        console.log('Models loaded:', models);
        displayModelsList(models);

        // Load first model if id in URL or if models available
        // Wait a moment for Three.js to fully initialize
        setTimeout(() => {
            const urlParams = new URLSearchParams(window.location.search);
            const id = urlParams.get('id');
            console.log('URL param id:', id);
            
            if (id) {
                const model = models.find(m => m.id == id);
                if (model) {
                    console.log('Found model from URL:', model);
                    loadModel(model);
                }
            } else if (models.length > 0) {
                console.log('Loading first model');
                loadModel(models[0]);
            } else {
                console.log('No models available');
                document.getElementById('loadingIndicator').textContent = 'No models available';
            }
        }, 100);
    } catch (error) {
        console.error('Error loading models:', error);
        document.getElementById('loadingIndicator').textContent = 'Error loading models: ' + error.message;
    }
}

function displayModelsList(modelsList) {
    const container = document.getElementById('modelList');
    container.innerHTML = modelsList.map(model => `
        <div class="model-item" onclick="selectModel(${model.id})">
            <h4>${model.name}</h4>
            <p>${(model.file_size / 1024).toFixed(2)} KB</p>
        </div>
    `).join('');
}

// Export selectModel for global use
window._selectModelFn = function(id) {
    const model = models.find(m => m.id === id);
    if (model) {
        loadModel(model);
        window.history.replaceState({}, '', `./viewer.html?id=${id}`);
    }
};

function loadModel(model) {
    console.log('Loading model:', model);
    document.getElementById('loadingIndicator').style.display = 'block';
    
    const loader = new GLTFLoader();
    
    // Remove existing model
    if (currentModel) {
        scene.remove(currentModel);
        currentModel = null;
    }

    const fileUrl = `http://localhost:8080${model.file_url}`;
    console.log('Requesting file from:', fileUrl);

    // Use fetch with Authorization header (supports archive tokens) then parse with GLTFLoader
    (async () => {
        try {
            const token = localStorage.getItem('token');
            const headers = {};
            if (token) headers['Authorization'] = `Bearer ${token}`;
            const resp = await fetch(fileUrl, { headers });
            if (!resp.ok) throw new Error('Failed to fetch model: ' + resp.statusText);
            const arrayBuffer = await resp.arrayBuffer();
            loader.parse(arrayBuffer, '', (gltf) => {
                console.log('Model loaded successfully via fetch+parse:', gltf);
                currentModel = gltf.scene;
            
            // Make sure model is visible
            currentModel.traverse((node) => {
                if (node.isMesh) {
                    node.castShadow = true;
                    node.receiveShadow = true;
                    console.log('Mesh found:', node.name);
                }
            });
            
            scene.add(currentModel);
            console.log('Model added to scene');

            // Auto-fit camera
            const box = new THREE.Box3().setFromObject(currentModel);
            const center = box.getCenter(new THREE.Vector3());
            const size = box.getSize(new THREE.Vector3());
            
            console.log('Model box:', { center, size });
            
            const maxDim = Math.max(size.x, size.y, size.z);
            const fov = camera.fov * (Math.PI / 180);
            let cameraZ = Math.abs(maxDim / 2 / Math.tan(fov / 2));
            cameraZ *= 1.5;

            // adjust camera clipping planes to avoid near-plane clipping when zooming
            // use a smaller near plane so users can zoom closer into components
            camera.near = Math.max(0.0001, maxDim / 10000);
            camera.far = Math.max(1000, cameraZ * 20);
            camera.updateProjectionMatrix();

            // position camera to frame the model
            camera.position.copy(center);
            camera.position.z += cameraZ;

            // adjust controls distances to sensible ranges based on model size
            // allow very close zooming by default; clamp to a small minimum
            controls.target.copy(center);
            controls.minDistance = Math.max(0.0001, maxDim * 0.001);
            controls.maxDistance = Math.max(cameraZ * 2, maxDim * 50);
            controls.screenSpacePanning = false;
            controls.update();

            console.log('Camera positioned at:', camera.position);

            // Update info
            document.getElementById('modelTitle').textContent = model.name;
            document.getElementById('modelDesc').textContent = model.description || 'No description';
            document.getElementById('modelInfo').textContent = `Uploaded by: ${model.uploaded_by} | Size: ${(model.file_size / 1024).toFixed(2)} KB`;
            document.getElementById('info-panel').style.display = 'block';
            document.getElementById('loadingIndicator').style.display = 'none';
            
            console.log('Model rendering complete');
            }, (error) => {
                console.error('Error parsing model:', error);
                document.getElementById('loadingIndicator').textContent = 'Error loading model: ' + (error.message || error);
            });
        } catch (err) {
            console.error('Error loading model via fetch:', err);
            document.getElementById('loadingIndicator').textContent = 'Error loading model: ' + err.message;
        }
    })();
}

// Initialize
checkAuth();
initThreeJS();
loadModelsList();
