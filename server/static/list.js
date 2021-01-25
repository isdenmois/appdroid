const loader = document.querySelector('.loader');
const list = document.querySelector('.apps-list');

export async function refetchList() {
  loader.classList.remove('hidden');
  list.classList.add('hidden');

  const apps = await fetch('/list').then((r) => r.json());

  list.innerHTML = '';
  for (const app of apps) {
    const item = document.createElement('li');
    item.innerHTML = `<div>${app.appId}</div><div>${app.version}</div>`;

    list.appendChild(item);
  }

  loader.classList.add('hidden');
  list.classList.remove('hidden');
}
