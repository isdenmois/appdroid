import { Injectable } from '@nestjs/common';
import { InjectModel } from '@nestjs/sequelize';
import { exec } from 'child_process';
import { renameSync } from 'fs';
import { Application } from './application.model';

@Injectable()
export class AppsService {
  constructor(
    @InjectModel(Application)
    private appModel: typeof Application,
  ) {}

  async addApk(path: string) {
    const appId = await getAppId(path);
    const version = await getAppVersion(path);

    await this.createOrUpdate(appId, version);

    renameSync(path, `static/${appId}.apk`);

    return { appId, version, url: `/static/${appId}.apk` };
  }

  async getList() {
    return this.appModel.findAll();
  }

  private async createOrUpdate(appId: string, version: string) {
    const model: Application = await this.appModel.findOne({
      where: { appId },
    });

    if (model === null) {
      this.appModel.create({ appId, version, type: 'static' });
    } else {
      model.update({ version });
    }
  }
}

function getAppId(path: string) {
  return apkanalyzer(`manifest application-id ${path}`);
}

function getAppVersion(path: string) {
  return apkanalyzer(`manifest version-code ${path}`);
}

function apkanalyzer(command: string) {
  return new Promise<string>((resolve, reject) =>
    exec(`apkanalyzer ${command}`, { encoding: 'utf8' }, (error, result) => {
      if (error) return reject(error);

      resolve(result.trim());
    }),
  );
}
