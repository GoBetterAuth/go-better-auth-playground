const ENV_CONFIG = {
  gobetterauth: {
    url: import.meta.env.VITE_GO_BETTER_AUTH_URL as string,
  },
} as const;

export default ENV_CONFIG;
