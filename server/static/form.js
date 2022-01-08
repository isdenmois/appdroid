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
    setUploading(true);

    if (target.files.length) {
      for (let i = 0; i < target.files.length; i++) {
        try {
          await uploadFile(target.files[i]);
        } catch (e) {
          console.error(e);
        }
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

  return fetch('/upload', {
    method: 'POST',
    body: formData,
  });
}
