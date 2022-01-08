export function AppList({ apps, onRemove }) {
  return apps.map((app) => AppPackage({ app, onRemove }));
}

function AppPackage({ app, onRemove }) {
  const item = document.createElement('li');
  item.innerHTML = `<a href="/${app.appId}.apk"><div class="name">${app.name}</div><div class="version">${app.versionName}</div></a>`;

  item.querySelector('.version').onclick = (e) => {
    e.preventDefault();

    removeApp(app, onRemove);
  };

  return item;
}

const removeApp = async (app, refetchList) => {
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
