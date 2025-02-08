import fs from 'node:fs';
import path from 'node:path';
import { tmpdir } from 'node:os';
import { randomUUID } from 'node:crypto';

const FILES_PATH = './data';
const TMP_PATH = path.join(tmpdir(), 'apks');

if (!fs.existsSync(TMP_PATH)) {
  fs.mkdirSync(TMP_PATH);
}

export const filesService = {
  async saveTemp(file: File) {
    const filePath = path.join(TMP_PATH, randomUUID());

    await Bun.write(filePath, file);

    return filePath;
  },

  async saveFile(from: string, name: string) {
    const to = path.join(FILES_PATH, name);

    await fs.promises.copyFile(from, to);

    return to;
  },

  getFile(apk: string) {
    const filePath = path.join(FILES_PATH, apk);

    return Bun.file(filePath);
  },

  async removeFile(apk: string) {
    const filePath = path.join(FILES_PATH, apk);

    await fs.promises.unlink(filePath);
  },
};
