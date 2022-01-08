const loader = document.querySelector('.loader');
const list = document.querySelector('.apps-list');

export async function refetchList() {
  loader.classList.remove('hidden');
  list.classList.add('hidden');

  const apps = await fetch('/list').then((r) => r.json());

  list.innerHTML = '';
  for (const app of apps) {
    const item = document.createElement('li');
    item.innerHTML = `<a href="/${app.appId}.apk">${app.appId}</a><div>${app.versionName}</div>`;

    item.querySelector('div').onclick = async () => {
      const toRemove = confirm(`Are you sure you want to remove "${app.name}"`);

      if (toRemove) {
        const response = await fetch(`/${app.appId}`, { method: 'DELETE' });

        if (response.ok) {
          refetchList();
        } else {
          alert("Can't remove the apk");
        }
      }
    };
    list.appendChild(item);
  }

  loader.classList.add('hidden');
  list.classList.remove('hidden');
}
