import {
  Controller,
  Delete,
  Get,
  Param,
  Post,
  UploadedFile,
  UseInterceptors,
} from '@nestjs/common';
import { FileInterceptor } from '@nestjs/platform-express';
import { ApiBody, ApiConsumes, ApiTags } from '@nestjs/swagger';
import { unlinkSync } from 'fs';
import { AppsService } from './apps.service';
import { FileUploadDto } from './dto/upload-dto';

@ApiTags('Apps')
@Controller()
export class AppsController {
  constructor(private appsService: AppsService) {}

  @Get('/list')
  getList() {
    return this.appsService.getList();
  }

  @Post('upload')
  @UseInterceptors(FileInterceptor('file'))
  @ApiConsumes('multipart/form-data')
  @ApiBody({
    description: 'Apk',
    type: FileUploadDto,
  })
  async uploadFile(@UploadedFile() file) {
    try {
      return await this.appsService.addApk(file.path);
    } catch (e) {
      unlinkSync(file.path);
    }
  }

  @Delete("/:appId")
  async removeFile(@Param("appId") appId: string) {
    return this.appsService.removeApk(appId)
  }
}
