import { withKey, isAuthed } from './auth.js';

const form = document.querySelector('form');
const fileInput = document.querySelector('input[type="file"]');

fileInput.addEventListener('focus', dragenter);
fileInput.addEventListener('dragenter', dragenter);
fileInput.addEventListener('dragleave', dragleave);
fileInput.addEventListener('blur', dragleave);
fileInput.addEventListener('drop', dragleave);

function dragenter() {
  form.classList.add('active');
}
function dragleave() {
  form.classList.remove('active');
}

export function UploadForm({ onUpload }) {
  async function onChange({ target }) {
    if (!isAuthed()) {
      alert('Enter your API key in the settings (⚙) before uploading.');
      return;
    }

    setUploading(true);

    if (target.files.length) {
      let failed = false;
      for (let i = 0; i < target.files.length; i++) {
        try {
          await uploadFile(target.files[i]);
        } catch (e) {
          failed = true;
          console.error(e);
        }
      }
      if (failed) {
        alert('Upload failed. Check your API key in the settings (⚙).');
      }
    }

    setUploading(false);
    target.value = '';
    onUpload();
  }

  fileInput.onchange = onChange;
}

const setUploading = (isVisible) =>
  form.classList.toggle('uploading', isVisible);

function uploadFile(file) {
  const formData = new FormData();
  formData.append('file', file);

  return fetch('/api/upload', withKey({ method: 'POST', body: formData })).then(
    async (response) => {
      if (!response.ok) {
        throw new Error(
          response.status === 401 ? 'Invalid API key' : 'Upload failed',
        );
      }
      return response;
    },
  );
}
