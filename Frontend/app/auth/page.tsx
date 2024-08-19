"use client"
import Image from 'next/image'
import React, { useEffect, useState } from 'react'
import Modal from '@/app/components/modal'
import Connectors from '../components/connectors'
import { useWeb3Modal } from '@web3modal/wagmi/react'
import { useAccount, useDisconnect } from 'wagmi'
import AuthAPI from '../utils/apis/auth'
import { signMessage } from 'wagmi/actions'
import { config } from '../web3/config'
import { cookie } from '../utils/cookie'
import { toast } from 'react-toastify'
import { useRouter } from 'next/navigation'
import ActivateUser from '../components/activateUser'
import { Link } from '../components/link'

const AuthHome = () => {
    const [showReg, setShowReg] = useState(false);
    const router = useRouter();
    const { open } = useWeb3Modal();
    const { disconnect } = useDisconnect();
    const { address, status } = useAccount();

    const { isPending: isVerifyingSignature, mutateAsync: verifySignature } = AuthAPI.verifySignature();
    const { isPending: isFetchingNounce, mutateAsync: fetchNounce, } = AuthAPI.useFetchNonce();
    const actions = (address: string) => new Promise(async (resolve, reject) => {
        await fetchNounce(address).then(async (res) => {
            const nonce = res.data.nonce;
            await signMessage(config, {
                message: nonce,
            }).then(async (signature) => {
                await verifySignature({ address: address!, signature }).then((credentialsResponse) => {
                    cookie.setJson('credentials', credentialsResponse.data.tokens);
                    resolve('done');
                    const { user } = credentialsResponse.data;
                    if (user.isActive) {
                        router.push("/u/h")
                    } else {
                        setShowReg(true);
                    }
                }).catch((e) => {
                    disconnect();
                    reject();
                })
            }).catch(() => {
                disconnect();
                reject();
            })
        }).catch((e) => {
            disconnect();
            reject();
        });
    });
    const handleConnectWallet = async () => {
        if (!isVerifyingSignature && !isFetchingNounce) {
            open();
        }
        else {
            toast.error("Something went wrong!")
        }
    }

    useEffect(() => {
        if (status === "connected") {
            (async () => {
                await toast.promise(actions(address!), {
                    error: "Something went wrong, please try again!",
                    pending: "Connecting wallet, please wait!!!",
                    success: "Wallet connected successfully!"
                });
            })();
        }
    }, [status])

    return (
        <main className='flex gap-8 items-center h-full'>
            {
                showReg && <Modal onTapOutside={() => setShowReg(false)}>
                    <ActivateUser />
                </Modal>
            }
            <Link href={"/i"} className='flex-[3] flex justify-center'>
                <Image src="/icon.png" alt="logo" width={200} height={200} />
            </Link>
            <article className='flex-[5]'>
                <h3 className='text-white text-6xl font-semibold max-w-lg'>Welcome to <span className='text-secondary'>PowerDfi</span> Studio</h3>
                <p className='text-white text-md mt-4 leading-relaxed max-w-lg'>The Studio allows Issuers to reserve, create and manage their security token.  connect wallet to get started</p>
                <button className='rounded-md bg-secondary w-full max-w-sm text-sm h-12 mt-8' onClick={handleConnectWallet}>Connect Wallet</button>
            </article>
        </main>
    )
}

export default AuthHome