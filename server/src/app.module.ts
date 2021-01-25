import { join } from 'path';
import { Module } from '@nestjs/common';
import { ServeStaticModule } from '@nestjs/serve-static';
import { SequelizeModule } from '@nestjs/sequelize';
import { AppsModule } from './apps/apps.module';

@Module({
  imports: [
    ServeStaticModule.forRoot({
      rootPath: join(__dirname, '..', 'static'),
    }),
    SequelizeModule.forRoot({
      dialect: 'sqlite',
      database: join(__dirname, '..', 'db.sqlite'),
      storage: join(__dirname, '..', 'db.sqlite'),
      autoLoadModels: true,
      synchronize: true,
    }),
    AppsModule,
  ],
})
export class AppModule {}
