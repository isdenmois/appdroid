// auth.js owns the API-key gate: a settings toggle that reveals a field where
// the key can be typed, sessionStorage-backed persistence, and helpers the
// upload/remove flows use to attach the X-API-Key header.
const STORAGE_KEY = 'appdroid-api-key';
const HEADER = 'X-API-Key';

// getKey returns the currently configured API key, or "" when none is set.
function getKey() {
  return sessionStorage.getItem(STORAGE_KEY) || '';
}

// setKey stores a key, clearing it when empty.
function setKey(value) {
  if (value) {
    sessionStorage.setItem(STORAGE_KEY, value);
  } else {
    sessionStorage.removeItem(STORAGE_KEY);
  }
}

// getApiKey exposes the current key to other modules.
export function getApiKey() {
  return getKey();
}

// isAuthed reports whether a key is present, i.e. whether mutating actions are
// unlocked.
export function isAuthed() {
  return getKey() !== '';
}

// withKey merges the X-API-Key header into a fetch init when a key is set.
export function withKey(init = {}) {
  const headers = { ...(init.headers || {}) };
  const key = getKey();
  if (key) {
    headers[HEADER] = key;
  }
  return { ...init, headers };
}

// setAuthed is the entry point called once the DOM is ready. It wires up the
// toggle, the field, and the lock/unlock indicator.
export function initAuth() {
  const toggle = document.getElementById('settings-toggle');
  const panel = document.getElementById('auth-panel');
  const input = document.getElementById('api-key-input');
  const saveBtn = document.getElementById('api-key-save');
  const clearBtn = document.getElementById('api-key-clear');
  const indicator = document.getElementById('auth-indicator');

  input.value = getKey();

  toggle.addEventListener('click', () => {
    panel.classList.toggle('hidden');
  });

  saveBtn.addEventListener('click', () => {
    setKey(input.value.trim());
    render(indicator);
    panel.classList.add('hidden');
  });

  clearBtn.addEventListener('click', () => {
    input.value = '';
    setKey('');
    render(indicator);
  });

  render(indicator);
}

// render reflects the auth state in the indicator: locked (no key) or unlocked.
function render(indicator) {
  const authed = isAuthed();
  indicator.classList.toggle('unlocked', authed);
  indicator.classList.toggle('locked', !authed);
  indicator.textContent = authed ? '🔓 unlocked' : '🔒 locked';
}
