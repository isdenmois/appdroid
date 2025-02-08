import aapt from 'aaptjs';

const PACKAGE_REGEXP =
  /name='([^']+)'\s*versionCode='(\d+)'\s*versionName='([^']+)/;
const NAME_REGEXP = /application:\s*label='([^']+)'/;

export const apkService = {
  async getInfo(path: string) {
    const info = await aapt.dump('badging', path);

    const [, appId, version, versionName] = info.match(PACKAGE_REGEXP) ?? [];
    const [, name] = info.match(NAME_REGEXP) ?? [];

    if (!appId || !name) {
      return null;
    }

    return {
      appId,
      name,
      version,
      versionName,
    };
  },
};
