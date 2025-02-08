import { Html } from '@elysiajs/html';

interface Props {
  title?: string;
  children: any;
}

export function Page({ title = '', children }: Props) {
  return (
    <html lang="en">
      <head>
        <title>{title}</title>

        <meta charset="UTF-8" />
        <meta
          name="viewport"
          content="width=device-width, initial-scale=1, maximum-scale=1, minimum-scale=1, user-scalable=no, minimal-ui"
        />
        <meta name="apple-mobile-web-app-capable" content="yes" />
        <meta name="apple-mobile-web-app-status-bar-style" content="black" />
        <link rel="stylesheet" href="/style.css" />
      </head>
      <body>{children}</body>
    </html>
  );
}
