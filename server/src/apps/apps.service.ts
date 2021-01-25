import { Injectable } from '@nestjs/common';
import { InjectModel } from '@nestjs/sequelize';
import { exec } from 'child_process';
import { renameSync } from 'fs';
import { Application } from './application.model';

type AppInfo = Pick<Application, 'appId' | 'version' | 'versionCode' | 'name'>;

@Injectable()
export class AppsService {
  constructor(
    @InjectModel(Application)
    private appModel: typeof Application,
  ) {}

  async addApk(path: string) {
    const info = await getInfo(path);

    await this.createOrUpdate(info);

    renameSync(path, `static/${info.appId}.apk`);

    return { ...info, url: `/static/${info.appId}.apk` };
  }

  async getList() {
    return this.appModel.findAll();
  }

  private async createOrUpdate(info: AppInfo) {
    const model: Application = await this.appModel.findOne({
      where: { appId: info.appId },
    });

    if (model === null) {
      this.appModel.create({ ...info, type: 'static' });
    } else {
      model.update(info);
    }
  }
}

const PACKAGE_REGEXP = /name='([^']+)'[\s]*versionCode='(\d+)'[\s]*versionName='([^']+)/;
const NAME_REGEXP = /:'(.*)'/;

async function getInfo(path: string): Promise<AppInfo> {
  const [, appId, versionCode, version] = (await aapt(path, 'package')).match(
    PACKAGE_REGEXP,
  );
  const [, name] = (await aapt(path, 'application-label')).match(NAME_REGEXP);

  return {
    appId,
    name,
    versionCode,
    version,
  };
}

function aapt(path, grep: string) {
  const cmd = ['aapt', 'dump', 'badging', path, '|', 'grep', grep].join(' ');

  return new Promise<string>((resolve, reject) =>
    exec(cmd, { encoding: 'utf8' }, (error, result) => {
      if (error) return reject(error);

      resolve(result.trim());
    }),
  );
}
