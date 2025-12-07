import { z } from 'zod';

export const oAuth2ProviderTypeSchema = z.enum(["discord", "github", "google"]);
export type OAuth2ProviderType = z.infer<typeof oAuth2ProviderTypeSchema>;

export const userSchema = z.object({
  id: z.uuid(),
  name: z.string().nonempty(),
  email: z.email(),
  emailVerified: z.boolean(),
  image: z.string().nullable().optional(),
  createdAt: z.string().nonempty(),
  updatedAt: z.string().nonempty(),
});
export type User = z.infer<typeof userSchema>;
