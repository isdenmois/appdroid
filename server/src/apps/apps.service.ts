import { Injectable, NotFoundException } from '@nestjs/common';
import { InjectModel } from '@nestjs/sequelize';
import { exec } from 'child_process';
import { renameSync, unlinkSync } from 'fs';
import { Application } from './application.model';

type AppInfo = Pick<Application, 'appId' | 'version' | 'versionName' | 'name'>;

@Injectable()
export class AppsService {
  constructor(
    @InjectModel(Application)
    private appModel: typeof Application,
  ) {}

  async addApk(path: string) {
    const info = await getInfo(path);

    await this.createOrUpdate(info);

    renameSync(path, `../static/${info.appId}.apk`);

    return { ...info, url: `/static/${info.appId}.apk` };
  }

  async getList() {
    return this.appModel.findAll();
  }

  async removeApk(appId: string) {
    const info = await this.appModel.findOne({ where: { appId } });

    if (info !== null) {
      await info.destroy();
      unlinkSync(`static/${info.appId}.apk`);
    } else {
      throw new NotFoundException();
    }
  }

  private async createOrUpdate(info: AppInfo) {
    const model: Application = await this.appModel.findOne({
      where: { appId: info.appId },
    });

    if (model === null) {
      await this.appModel.create({ ...info, type: 'static' } as Application);
    } else {
      await model.update(info);
    }
  }
}

const PACKAGE_REGEXP =
  /name='([^']+)'[\s]*versionCode='(\d+)'[\s]*versionName='([^']+)/;
const NAME_REGEXP = /:'(.*)'/;

async function getInfo(path: string): Promise<AppInfo> {
  const [, appId, version, versionName] = (await aapt(path, 'package')).match(
    PACKAGE_REGEXP,
  );
  const [, name] = (await aapt(path, 'application-label')).match(NAME_REGEXP);

  return {
    appId,
    name,
    version,
    versionName,
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
