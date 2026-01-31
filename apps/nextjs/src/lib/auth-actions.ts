"use server";

import { toCamelCaseKeys } from "es-toolkit";
import { z } from "zod";

import { sessionSchema, userSchema } from "@/models";
import { wrappedFetch } from "./gobetterauth-server";
import { applySetCookieHeaders } from "./gobetterauth-server-actions";

/**
 * Server Action to get current user and apply CSRF token cookie
 * This should be called from client components or form actions
 */
export async function getMeAction() {
  const result = await wrappedFetch("/me", {
    method: "GET",
    includeCookies: true,
    throwOnError: false,
  });

  // Apply cookies in Server Action context
  await applySetCookieHeaders(result.setCookieHeaders);

  if (result.error) {
    return null;
  }

  const validationResult = z
    .object({
      user: userSchema,
      session: sessionSchema,
    })
    .parse(toCamelCaseKeys(result.data));

  return validationResult;
}
