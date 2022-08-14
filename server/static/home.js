import { AppList } from './list.js';
import { UploadForm } from './form.js';

const loader = document.querySelector('.loader');
const list = document.querySelector('.apps-list');

async function refetchList() {
  setLoader(true);

  const apps = await fetch('/api/list').then((r) => r.json());

  setApps(sortBy(apps, 'name'));
  setLoader(false);
}

const setLoader = (isVisible) => {
  loader.classList.toggle('hidden', !isVisible);
  list.classList.toggle('hidden', isVisible);
};

const setApps = (apps) => {
  list.innerHTML = '';

  for (const item of AppList({ apps, onRemove: refetchList })) {
    list.appendChild(item);
  }
};

export const Home = () => {
  refetchList();

  UploadForm({ onUpload: refetchList });
  AppList({ apps: [] });
};

const sortBy = (array, key) =>
  array.sort((a, b) => {
    if (a[key] < b[key]) {
      return -1;
    }
    if (a[key] > b[key]) {
      return 1;
    }

    return 0;
  });
