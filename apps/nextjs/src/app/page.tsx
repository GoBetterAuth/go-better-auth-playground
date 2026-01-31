import Image from 'next/image';
import Link from 'next/link';

import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';

export default function HomePage() {
  return (
    <main className="flex min-h-screen flex-col items-center justify-center p-24">
      <Card className="w-full max-w-md text-center">
        <CardHeader className="flex flex-col items-center">
          <Image
            className="relative dark:drop-shadow-[0_0_0.3rem_#ffffff70] dark:invert"
            src="/app-logo.png"
            alt="GoBetterAuth Logo"
            width={180}
            height={37}
            priority
          />
          <CardTitle className="mt-4 text-2xl font-semibold">
            Welcome to GoBetterAuth Playground
          </CardTitle>
          <CardDescription>
            An example modern authentication solution built with Go and Next.js.
          </CardDescription>
        </CardHeader>
        <CardContent className="flex flex-col gap-4">
          <Link href="/auth/sign-in" passHref>
            <Button className="w-full">Sign In</Button>
          </Link>
          <Link href="/auth/sign-up" passHref>
            <Button className="w-full" variant="outline">
              Sign Up
            </Button>
          </Link>
        </CardContent>
      </Card>
    </main>
  );
}
