import path from 'node:path';
import fs from 'node:fs/promises';
import os from 'node:os';
import { exec } from 'node:child_process';

const platform = os.platform();
const aapt = path.join('node_modules', 'aaptjs', 'bin', platform, 'aapt');

if (platform === 'linux') {
  fs.chmod(aapt, '777');
}

export const aaptService = {
  dump(path: string) {
    const cmd = [aapt, 'dump', 'badging', path].join(' ');

    return new Promise<string>((resolve, reject) =>
      exec(cmd, { encoding: 'utf8' }, (error, result) => {
        if (error) return reject(error);

        resolve(result.trim());
      }),
    );
  },
};
