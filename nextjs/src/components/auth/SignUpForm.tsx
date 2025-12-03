"use client";

import { useAction } from 'next-safe-action/hooks';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { toast } from 'sonner';
import { z } from 'zod';

import { signUpAction } from '@/app/actions';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { useForm } from '@tanstack/react-form';

const formSchema = z
  .object({
    name: z.string().nonempty("Name is required"),
    email: z.email("Invalid email address"),
    password: z.string().nonempty("Password is required"),
    confirmPassword: z.string().nonempty("Confirm Password is required"),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: "Passwords don't match",
    path: ["confirmPassword"],
  });

export default function SignUpForm() {
  const router = useRouter();

  const { executeAsync } = useAction(signUpAction);

  const form = useForm({
    defaultValues: {
      name: "John Doe",
      email: "john.doe@example.com",
      password: "Pass!2345",
      confirmPassword: "Pass!2345",
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
        toast.success("Signed up successfully!");
        router.push("/auth/email-verification");
      } catch (error: any) {
        console.error("Error during registration:", error);
        toast.error(error.message);
      }
    },
  });

  return (
    <Card className="w-full max-w-md">
      <CardHeader className="text-center">
        <CardTitle className="text-2xl font-bold">Sign Up</CardTitle>
        <CardDescription>Create your account to get started.</CardDescription>
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
            name="name"
            children={(field) => (
              <div className="grid gap-2">
                <Label htmlFor={field.name}>Name</Label>
                <Input
                  id={field.name}
                  type="text"
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

          <form.Field
            name="confirmPassword"
            children={(field) => (
              <div className="grid gap-2">
                <Label htmlFor={field.name}>Confirm Password</Label>
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
                {isSubmitting ? "Registering..." : "Sign Up"}
              </Button>
            )}
          />
          <div className="mt-4 text-center text-sm">
            Already have an account?{" "}
            <Link href="/auth/sign-in" className="underline">
              Sign In
            </Link>
          </div>
        </form>
      </CardContent>
    </Card>
  );
}
