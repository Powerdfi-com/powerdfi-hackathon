/* eslint-disable react-hooks/exhaustive-deps */
"use client"
import React, { ReactNode, useCallback, useEffect, useState } from 'react'
import UserAPI from '../utils/apis/user';
import { useAccount, useDisconnect } from 'wagmi';
import Loading from '../components/loading';
import { TCredentials, TUser } from '../utils/types';
import { dummyUser, UserContext } from '../context/context';
import { useRouter } from 'next/navigation';
import { toast } from 'react-toastify';
import { cookie } from '../utils/cookie';
import AuthAPI from '../utils/apis/auth';

const ContextLayout = ({ children }: { children: ReactNode }) => {
    const [credentials, setCredentials] = useState<TCredentials>(cookie.getJson("credentials"));
    const { isIdle, isPending, mutateAsync: fetchUser } = UserAPI.getProfile();
    const { data, isPending: isPendingWallet } = UserAPI.getWalletDetails();
    const router = useRouter()
    const { isConnected } = useAccount();
    const { disconnect } = useDisconnect();
    const signOut = useCallback(() => {
        toast.error("Credentials expired, please sign in again!")
        cookie.remove('credentials');
        disconnect();
        router.push("/auth/")
    }, [disconnect, router])
    const { mutateAsync: fetchRefreshToken } = AuthAPI.refreshToken({
        onSuccess: (res) => {
            cookie.setJson("credentials", res.data);
            setCredentials(res.data);
        },
        onError: (error) => {
            // Refresh token has expired...
            signOut();
        }
    })
    const [user, setUser] = useState<TUser>(dummyUser);

    const handleClickCopyAddress = async () => {
        await navigator.clipboard.writeText(data?.data.address!).then(() => {
            toast.success("Address copied to clipboard!")
        }).catch((e) => {
            toast.error("Clipboard not supported, please check permission and retry!")
        })
    }

    useEffect(() => {
        const handleLoadFetchUSer = async () => {
            await fetchUser().then((res) => {
                setUser(res.data);
            }).catch(() => {
                toast.error("Error fetching user... Trying again!")
            })
        }
        if (credentials && isConnected) {
            if (new Date(credentials.expiresAt * 1000) > new Date()) {
                handleLoadFetchUSer();
            }
            else {
                fetchRefreshToken(credentials.refreshToken);
            }
        } else {
            signOut();
        }
    }, [])

    // Fix login issue in case wallet is disconnected on wallet itself

    if (!isConnected) return <div></div>;

    return (isPending || isPendingWallet || isIdle || !user.email) ? <Loading className="!h-screen !w-screen" /> : <UserContext.Provider value={{ signOut, user, setUser, address: data!.data.address, copyAddress: handleClickCopyAddress }}>
        {children}
    </UserContext.Provider>
}

export default ContextLayout