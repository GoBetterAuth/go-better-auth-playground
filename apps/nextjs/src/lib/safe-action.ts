import {
  createSafeActionClient,
  DEFAULT_SERVER_ERROR_MESSAGE,
} from "next-safe-action";
import * as z from "zod";

import { ActionError } from "@/models";
import { goBetterAuthClientServer } from "./gba-client-server";
import { GetMeResponse } from "go-better-auth";

export const actionClient = createSafeActionClient({
  defineMetadataSchema() {
    return z.object({
      actionName: z.string(),
    });
  },
  handleServerError(e) {
    console.error("Action error:", e.message);

    if (e instanceof ActionError) {
      return e.message;
    }

    return DEFAULT_SERVER_ERROR_MESSAGE;
  },
});

// Auth client defined by extending the base one.
// Note that the same initialization options and middleware functions of the base client
// will also be used for this one.
export const authActionClient = actionClient.use(async ({ next }) => {
  try {
    const data = await goBetterAuthClientServer.getMe<GetMeResponse>();
    if (!data.user) {
      throw new Error("Unauthorized");
    }

    return next({ ctx: { userId: data.user.id } });
  } catch (error) {
    throw new Error(error instanceof Error ? error.message : String(error));
  }
});
