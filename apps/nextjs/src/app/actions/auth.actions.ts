"use server";

import { redirect } from "next/navigation";
import { flattenValidationErrors } from "next-safe-action";

import { z } from "zod";

import { ENV_CONFIG } from "@/constants/env-config";
import { actionClient } from "@/lib/safe-action";
import { ActionError } from "@/models";
import { goBetterAuthClientServer } from "@/lib/gba-client-server";
import { GetMeResponse } from "go-better-auth";

// --------------------------

const signUpFormSchema = z.object({
  name: z.string().nonempty("Name is required"),
  email: z.email("Invalid email address"),
  password: z
    .string()
    .min(8, "Password is required")
    .max(32, "Password must be at most 32 characters"),
});

export const signUpAction = actionClient
  .metadata({ actionName: "signUpAction" })
  .inputSchema(signUpFormSchema, {
    handleValidationErrorsShape: async (ve, utils) =>
      flattenValidationErrors(ve).fieldErrors,
  })
  .action(async ({ parsedInput }) => {
    try {
      const data: any = await goBetterAuthClientServer.emailPassword.signUp({
        name: parsedInput.name,
        email: parsedInput.email,
        password: parsedInput.password,
        callbackUrl: `${ENV_CONFIG.baseUrl}/dashboard`,
      });

      return data;
    } catch (error) {
      console.error("Error signing up:", error);
      if (error instanceof ActionError) {
        throw error;
      }
      throw new ActionError(
        error instanceof Error
          ? error.message
          : "Failed to sign up user. Please try again later.",
      );
    }
  });

// --------------------------

const signInFormSchema = z.object({
  email: z.email("Invalid email address"),
  password: z.string().nonempty(),
});

export const signInAction = actionClient
  .metadata({ actionName: "signInAction" })
  .inputSchema(signInFormSchema, {
    handleValidationErrorsShape: async (ve, utils) =>
      flattenValidationErrors(ve).fieldErrors,
  })
  .action(async ({ parsedInput }) => {
    let data: GetMeResponse | null = null;
    try {
      data = await goBetterAuthClientServer.emailPassword.signIn({
        email: parsedInput.email,
        password: parsedInput.password,
        callbackUrl: `${ENV_CONFIG.baseUrl}/dashboard`,
      });
      console.log(data);
    } catch (error) {
      console.error("Error signing in:", error);
      if (error instanceof ActionError) {
        throw error;
      }
      throw new ActionError(
        error instanceof Error
          ? error.message
          : "Failed to sign in user. Please try again later.",
      );
    }

    if (data && !data.user.emailVerified) {
      return redirect(`/auth/email-verification?email=${data.user.email}`);
    }

    return redirect("/dashboard");
  });

// --------------------------

export const sendEmailVerificationAction = actionClient
  .metadata({ actionName: "sendEmailVerificationAction" })
  .inputSchema(
    z.object({
      email: z.email("Invalid email address"),
    }),
  )
  .action(async ({ parsedInput }) => {
    try {
      const data: any =
        await goBetterAuthClientServer.emailPassword.sendEmailVerification({
          email: parsedInput.email,
          callbackUrl: `${ENV_CONFIG.baseUrl}/dashboard`,
        });

      return data;
    } catch (error) {
      console.error("Error sending email verification:", error);
      if (error instanceof ActionError) {
        throw error;
      }
      throw new ActionError(
        error instanceof Error
          ? error.message
          : "Failed to send email verification. Please try again later.",
      );
    }
  });

// --------------------------

const resetPasswordFormSchema = z.object({
  email: z.email("Invalid email address"),
});

export const resetPasswordAction = actionClient
  .metadata({ actionName: "resetPasswordAction" })
  .inputSchema(resetPasswordFormSchema, {
    handleValidationErrorsShape: async (ve, utils) =>
      flattenValidationErrors(ve).fieldErrors,
  })
  .action(async ({ parsedInput }) => {
    try {
      const data: any =
        await goBetterAuthClientServer.emailPassword.requestPasswordReset({
          email: parsedInput.email,
          callbackUrl: `${ENV_CONFIG.baseUrl}/auth/change-password`,
        });

      return data;
    } catch (error) {
      console.error("Error signing in:", error);
      if (error instanceof ActionError) {
        throw error;
      }
      throw new ActionError(
        error instanceof Error
          ? error.message
          : "Failed to send reset password link. Please try again later.",
      );
    }
  });

// --------------------------

const changePasswordFormSchema = z.object({
  token: z.string().nonempty("Token is required"),
  newPassword: z
    .string()
    .min(8, "Password is required")
    .max(32, "Password must be at most 32 characters"),
});

export const changePasswordAction = actionClient
  .metadata({ actionName: "changePasswordAction" })
  .inputSchema(changePasswordFormSchema, {
    handleValidationErrorsShape: async (ve, utils) =>
      flattenValidationErrors(ve).fieldErrors,
  })
  .action(async ({ parsedInput }) => {
    try {
      const data: any =
        await goBetterAuthClientServer.emailPassword.changePassword({
          token: parsedInput.token,
          password: parsedInput.newPassword,
        });

      return data;
    } catch (error) {
      console.error("Error changing password:", error);
      if (error instanceof ActionError) {
        throw error;
      }
      throw new ActionError(
        error instanceof Error
          ? error.message
          : "Failed to change password. Please try again later.",
      );
    }
  });

// --------------------------

const emailChangeFormSchema = z
  .object({
    newEmail: z.email("Invalid email address").trim(),
    confirmEmail: z.email("Invalid email address").trim(),
  })
  .refine((data) => data.newEmail === data.confirmEmail, {
    message: "Emails don't match",
    path: ["confirmEmail"],
  });

export const emailChangeAction = actionClient
  .metadata({ actionName: "emailChangeAction" })
  .inputSchema(emailChangeFormSchema, {
    handleValidationErrorsShape: async (ve, utils) =>
      flattenValidationErrors(ve).fieldErrors,
  })
  .action(async ({ parsedInput }) => {
    try {
      const data: any =
        await goBetterAuthClientServer.emailPassword.requestEmailChange({
          email: parsedInput.newEmail,
          callbackUrl: `${ENV_CONFIG.baseUrl}/dashboard`,
        });

      return data;
    } catch (error) {
      console.error("Error changing email:", error);
      if (error instanceof ActionError) {
        throw error;
      }
      throw new ActionError(
        error instanceof Error
          ? error.message
          : "Failed to change email. Please try again later.",
      );
    }
  });

// --------------------------
