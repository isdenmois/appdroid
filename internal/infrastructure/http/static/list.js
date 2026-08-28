import { withKey, isAuthed } from './auth.js';

export function AppList({ apps, onRemove }) {
  return apps.map((app) => AppPackage({ app, onRemove }));
}

function AppPackage({ app, onRemove }) {
  const item = document.createElement('li');
  item.innerHTML = `<a href="/file/${app.apk}"><div class="name">${app.name}</div><div class="version">${app.versionName}</div></a>`;

  item.querySelector('.version').onclick = (e) => {
    e.preventDefault();

    removeApp(app, onRemove);
  };

  return item;
}

const removeApp = async (app, refetchList) => {
  if (!isAuthed()) {
    alert('Enter your API key in the settings (⚙) before removing.');
    return;
  }

  const toRemove = confirm(`Are you sure you want to remove "${app.name}"`);

  if (toRemove) {
    const response = await fetch(
      `/api/${app.appId}`,
      withKey({ method: 'DELETE' }),
    );

    if (response.ok) {
      refetchList();
    } else if (response.status === 401) {
      alert("Invalid API key — can't remove the apk.");
    } else {
      alert("Can't remove the apk");
    }
  }
};
