import Elysia from 'elysia';

export const pingController = new Elysia().get('/ping', () => 'pong');
