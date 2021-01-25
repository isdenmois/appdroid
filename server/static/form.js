import { refetchList } from './list.js';

const form = document.querySelector('form');
const fileInput = document.querySelector('input[type="file"]');

fileInput.addEventListener('change', onChange);
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

async function onChange({ target }) {
  form.classList.add('uploading');
  if (target.files.length) {
    for (let i = 0; i < target.files.length; i++) {
      try {
        await uploadFile(target.files[i]);
      } catch (e) {
        console.error(e);
      }
    }
  }

  form.classList.remove('uploading');
  target.value = '';
  refetchList();
}

function uploadFile(file) {
  const formData = new FormData();
  formData.append('file', file);

  return fetch('/upload', {
    method: 'POST',
    body: formData,
  });
}
