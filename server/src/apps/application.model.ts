import { Column, Model, Table } from 'sequelize-typescript';

@Table
export class Application extends Model<Application> {
  @Column
  name: string;

  @Column
  appId: string;

  @Column
  version: string;

  @Column
  type: string;
}
