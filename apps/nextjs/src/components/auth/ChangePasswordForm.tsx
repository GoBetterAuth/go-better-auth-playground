"use client";

import { useAction } from "next-safe-action/hooks";
import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { toast } from "sonner";
import { z } from "zod";

import { changePasswordAction } from "@/app/actions";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useForm } from "@tanstack/react-form";

const formSchema = z.object({
  newPassword: z
    .string()
    .min(8, "Password must be at least 8 characters")
    .max(32, "Password must be at most 32 characters"),
});

export default function ChangePasswordForm() {
  const searchParams = useSearchParams();
  const token = searchParams.get("token") || "";
  const router = useRouter();

  const { executeAsync } = useAction(changePasswordAction);

  const form = useForm({
    defaultValues: {
      newPassword: "",
    },
    validators: {
      onChange: formSchema,
    },
    onSubmit: async ({ value }) => {
      try {
        const data = await executeAsync({
          token,
          newPassword: value.newPassword,
        });
        if (data.serverError) {
          throw new Error(data.serverError);
        }
        if (data.validationErrors) {
          throw new Error(
            Object.entries(data.validationErrors)
              .map(([_, v]) => v)
              .join(", "),
          );
        }
        console.log(data);
        toast.success("Changed password successfully!");
        router.push("/dashboard");
      } catch (error: any) {
        console.error("Error during password change:", error);
        toast.error(error.message);
      }
    },
  });

  return (
    <Card className="w-full max-w-md">
      <CardHeader className="text-center">
        <CardTitle className="text-2xl font-bold">Change Password</CardTitle>
        <CardDescription>
          Enter your new password to update your account.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form
          onSubmit={(e) => {
            e.preventDefault();
            e.stopPropagation();
            form.handleSubmit();
          }}
          className="grid gap-4"
        >
          <form.Field
            name="newPassword"
            children={(field) => (
              <div className="grid gap-2">
                <Label htmlFor={field.name}>New Password</Label>
                <Input
                  id={field.name}
                  type="password"
                  value={field.state.value}
                  onBlur={field.handleBlur}
                  onChange={(e) => field.handleChange(e.target.value)}
                />
                {field.state.meta.isTouched && !field.state.meta.isValid ? (
                  <em className="text-red-500 text-sm">
                    {field.state.meta.errors
                      .map((error) => error?.message)
                      .join(", ")}
                  </em>
                ) : null}
              </div>
            )}
          />

          <form.Subscribe
            selector={(state) => [state.canSubmit, state.isSubmitting]}
            children={([canSubmit, isSubmitting]) => (
              <Button type="submit" className="w-full" disabled={!canSubmit}>
                {isSubmitting ? "Changing password..." : "Change Password"}
              </Button>
            )}
          />
          <div className="mt-4 text-center text-sm">
            Want to sign in?{" "}
            <Link href="/auth/sign-in" className="underline">
              Sign in
            </Link>
          </div>
        </form>
      </CardContent>
    </Card>
  );
}
