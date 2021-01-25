import { Module } from '@nestjs/common';
import { SequelizeModule } from '@nestjs/sequelize';
import { MulterModule } from '@nestjs/platform-express';
import * as multer from 'multer';
import { Application } from './application.model';
import { AppsController } from './apps.controller';
import { AppsService } from './apps.service';

const storage = multer.diskStorage({
  destination: 'upload',
});

@Module({
  imports: [
    SequelizeModule.forFeature([Application]),
    MulterModule.register({
      dest: 'upload',
    }),
  ],
  controllers: [AppsController],
  providers: [AppsService],
})
export class AppsModule {}
