import { createClient } from "go-better-auth";
import {
  EmailPasswordPlugin,
  OAuth2Plugin,
  CSRFPlugin,
  JWTPlugin,
  BearerPlugin,
} from "go-better-auth/plugins";

import ENV_CONFIG from "~/constants/env-config";

export const goBetterAuthClient = createClient({
  url: ENV_CONFIG.gobetterauth.url,
  plugins: [
    new EmailPasswordPlugin(),
    new OAuth2Plugin(),
    new CSRFPlugin({
      cookieName: "gobetterauth_csrf_token",
      headerName: "x-gobetterauth-csrf-token",
    }),
    new JWTPlugin(),
    new BearerPlugin(),
  ],
});
