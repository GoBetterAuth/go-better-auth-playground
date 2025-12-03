"use client";

import { useAction } from 'next-safe-action/hooks';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { toast } from 'sonner';
import { z } from 'zod';

import { signInAction } from '@/app/actions';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { useForm } from '@tanstack/react-form';

const formSchema = z.object({
  email: z.email("Invalid email address"),
  password: z.string().nonempty("Password is required"),
});

export default function SignInForm() {
  const router = useRouter();

  const { executeAsync } = useAction(signInAction);

  const form = useForm({
    defaultValues: {
      email: "john.doe@example.com",
      password: "Pass!2345",
    },
    validators: {
      onChange: formSchema,
    },
    onSubmit: async ({ value }) => {
      try {
        const data = await executeAsync(value);
        if (data.serverError) {
          throw new Error(data.serverError);
        }
        if (data.validationErrors) {
          throw new Error(
            Object.entries(data.validationErrors)
              .map(([_, v]) => v)
              .join(", ")
          );
        }
        console.log(data);
        toast.success("Signed in successfully!");
        router.push("/dashboard");
      } catch (error: any) {
        console.error("Error during sign in:", error);
        toast.error(error.message);
      }
    },
  });

  return (
    <Card className="w-full max-w-md">
      <CardHeader className="text-center">
        <CardTitle className="text-2xl font-bold">Sign In</CardTitle>
        <CardDescription>
          Enter your email and password to access your account.
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
            name="email"
            children={(field) => (
              <div className="grid gap-2">
                <Label htmlFor={field.name}>Email</Label>
                <Input
                  id={field.name}
                  type="email"
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

          <form.Field
            name="password"
            children={(field) => (
              <div className="grid gap-2">
                <Label htmlFor={field.name}>Password</Label>
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

          <div className="text-right text-sm">
            <Link href="/auth/reset-password" className="underline">
              Forgot password?
            </Link>
          </div>

          <form.Subscribe
            selector={(state) => [state.canSubmit, state.isSubmitting]}
            children={([canSubmit, isSubmitting]) => (
              <Button type="submit" className="w-full" disabled={!canSubmit}>
                {isSubmitting ? "Signing in..." : "Sign In"}
              </Button>
            )}
          />
          <div className="mt-4 text-center text-sm">
            Don&apos;t have an account?{" "}
            <Link href="/auth/sign-up" className="underline">
              Sign up
            </Link>
          </div>
        </form>
      </CardContent>
    </Card>
  );
}
