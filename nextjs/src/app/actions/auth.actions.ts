"use server";

import { flattenValidationErrors } from 'next-safe-action';

import { z } from 'zod';

import { goBetterAuthServer } from '@/lib/gobetterauth-server';
import { actionClient } from '@/lib/safe-action';
import { ActionError } from '@/models';

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
      const data: { token: string; user: any } =
        await goBetterAuthServer.signUp.email(
          parsedInput.name,
          parsedInput.email,
          parsedInput.password,
          "http://localhost:3000/dashboard"
        );

      return data;
    } catch (error) {
      console.error("Error signing up:", error);
      if (error instanceof ActionError) {
        throw error;
      }
      throw new ActionError(
        error instanceof Error
          ? error.message
          : "Failed to sign up user. Please try again later."
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
    try {
      const data: { token: string; user: any } =
        await goBetterAuthServer.signIn.email(
          parsedInput.email,
          parsedInput.password,
          "http://localhost:3000/dashboard"
        );

      return data;
    } catch (error) {
      console.error("Error signing in:", error);
      if (error instanceof ActionError) {
        throw error;
      }
      throw new ActionError(
        error instanceof Error
          ? error.message
          : "Failed to sign in user. Please try again later."
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
      const data = await goBetterAuthServer.resetPassword(
        parsedInput.email,
        "http://localhost:3000/auth/change-password"
      );

      return data;
    } catch (error) {
      console.error("Error signing in:", error);
      if (error instanceof ActionError) {
        throw error;
      }
      throw new ActionError(
        error instanceof Error
          ? error.message
          : "Failed to send reset password link. Please try again later."
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
      const data = await goBetterAuthServer.changePassword(
        parsedInput.token,
        parsedInput.newPassword
      );

      return data;
    } catch (error) {
      console.error("Error signing in:", error);
      if (error instanceof ActionError) {
        throw error;
      }
      throw new ActionError(
        error instanceof Error
          ? error.message
          : "Failed to change password. Please try again later."
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
      const data = await goBetterAuthServer.emailChange(
        parsedInput.newEmail,
        "http://localhost:3000/dashboard"
      );

      return data;
    } catch (error) {
      console.error("Error changing email:", error);
      if (error instanceof ActionError) {
        throw error;
      }
      throw new ActionError(
        error instanceof Error
          ? error.message
          : "Failed to change email. Please try again later."
      );
    }
  });

// --------------------------

export const signOutAction = actionClient
  .metadata({ actionName: "signOutAction" })
  .action(async () => {
    try {
      const data = await goBetterAuthServer.signOut();

      return data;
    } catch (error) {
      console.error("Error signing out:", error);
      if (error instanceof ActionError) {
        throw error;
      }
      throw new ActionError(
        error instanceof Error
          ? error.message
          : "Failed to sign out. Please try again later."
      );
    }
  });
