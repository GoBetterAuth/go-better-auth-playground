import { NextRequest, NextResponse } from 'next/server';

import { goBetterAuthServer } from './lib/gobetterauth-server';

export async function proxy(request: NextRequest) {
  try {
    const data = await goBetterAuthServer.getSession();

    const isAuthenticated = !!data?.user;

    if (isAuthenticated) {
      if (!data.user.emailVerified) {
        if (!request.nextUrl.pathname.startsWith("/auth/email-verification")) {
          const url = request.nextUrl.clone();
          url.pathname = "/auth/email-verification";
          return NextResponse.redirect(url);
        }
        // Allow access to email verification page
        return NextResponse.next();
      }
      // Only redirect to dashboard if not already there and not on email verification page
      if (
        !request.nextUrl.pathname.startsWith("/dashboard") &&
        !request.nextUrl.pathname.startsWith("/auth/email-verification")
      ) {
        const url = request.nextUrl.clone();
        url.pathname = "/dashboard";
        return NextResponse.redirect(url);
      }
    } else {
      // Only redirect if trying to access dashboard
      if (request.nextUrl.pathname.startsWith("/dashboard")) {
        const url = request.nextUrl.clone();
        url.pathname = "/auth/sign-in";
        return NextResponse.redirect(url);
      }
    }
  } catch (error: any) {
    if (!request.nextUrl.pathname.startsWith("/auth/sign-in")) {
      const url = request.nextUrl.clone();
      url.pathname = "/auth/sign-in";
      return NextResponse.redirect(url);
    }
  }

  return NextResponse.next();
}

export const config = {
  matcher: [
    /*
     * Match all request paths except for the ones starting with:
     * - _next/static (static files)
     * - _next/image (image optimization files)
     * - favicon.ico (favicon file)
     * Feel free to modify this pattern to include more paths.
     */
    "/((?!_next/static|_next/image|favicon.ico|api|.*\\.(?:svg|png|jpg|jpeg|gif|webp)$).*)",
  ],
};
