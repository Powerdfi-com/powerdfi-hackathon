'use client';

import { useTransition } from 'react';
import NextLink from 'next/link';
import NProgress from 'nprogress';
import { useEffect } from "react";
import { useRouter } from 'next/navigation';

/**
 * A custom Link component that wraps Next.js's next/link component.
 */
export function Link({
    href,
    children,
    replace,
    ...rest
}: Parameters<typeof NextLink>[0]) {
    const router = useRouter();
    const [isPending, startTransition] = useTransition();
    useEffect(() => {
        if (isPending) {
            NProgress.start();
        }
        else {
            NProgress.done();
        }
        return () => {
            NProgress.done()
        };
    }, [isPending])

    return (
        <NextLink
            href={href}
            onClick={(e) => {
                e.preventDefault();
                startTransition(() => {
                    const url = href.toString();
                    if (replace) {
                        router.replace(url);
                    } else {
                        router.push(url);
                    }
                });
            }}
            {...rest}
        >
            {children}
        </NextLink>
    );
}