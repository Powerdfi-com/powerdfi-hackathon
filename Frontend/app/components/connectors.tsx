"use client"
import Image from 'next/image'
import React, { useState } from 'react'
import { Connector, useConnect, useDisconnect } from 'wagmi'
import AuthAPI from '../utils/apis/auth'
import { signMessage } from 'wagmi/actions'
import { config } from '../web3/config'
import { cookie } from '../utils/cookie'
import { useRouter } from 'next/navigation'
import { toast } from 'react-toastify'
import Modal from './modal'
import ActivateUser from './activateUser'

const Connectors = () => {
    const tabs = [
        "Etherum", "Polygon", "BSC", "Solana"
    ]
    const [tab, setTab] = useState("Etherum");
    const { connectors, connect, } = useConnect();
    const router = useRouter();
    const [showReg, setShowReg] = useState(false);
    const { disconnect } = useDisconnect();

    const { isPending: isVerifyingSignature, mutateAsync: verifySignature } = AuthAPI.verifySignature();
    const { isPending: isFetchingNounce, mutateAsync: fetchNounce, } = AuthAPI.useFetchNonce();


    const handleConnectWallet = async (connector: Connector) => {
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
        if (!isVerifyingSignature && !isFetchingNounce) {
            await connector.connect().then(async (res) => {
                const address = res.accounts[0];
                await toast.promise(actions(address), {
                    error: "Something went wrong, please try again!",
                    pending: "Connecting wallet, please wait!!!",
                    success: "Wallet connected successfully!"
                });
            }).catch((e) => {
                toast.error("Wallet not connected!")
            })
        }
        else {
            toast.error("Something went wrong!")
        }
    }
    return <>
        {
            showReg && <Modal onTapOutside={() => setShowReg(false)}>
                <ActivateUser />
            </Modal>
        }
        <div className='ring-1 w-full max-w-md ring-primary/40 rounded-xl bg-black-shade py-12 px-8 sm:px-12 md:px-24 flex flex-col items-center gap-6' onClick={(e) => e.stopPropagation()}>
            <div className='h-16 w-16 relative'>
                <Image src={"/logo.png"} fill={true} alt="logo" className='object-cover' />
            </div>
            <article>
                <h4 className="text-white font-semibold text-2xl text-center">Connect Wallet</h4>
                <p className='text-white/80 text-xs text-center'>Choose a wallet to connect</p>
            </article>
            <section className='flex flex-col gap-6 w-full'>
                <ul className='flex justify-evenly'>
                    {
                        tabs.map((_tab) => <li key={_tab}><div className='border-b border-secondary text-white text-xs py-0.5'>{_tab}</div></li>)
                    }
                </ul>
                <ul className='flex flex-col gap-4'>
                    {connectors.map((connector) => <li key={connector.id}>
                        <div onClick={() => handleConnectWallet(connector)} className='ring-1 ring-secondary rounded-xl flex items-center gap-4 w-full cursor-pointer'>
                            <div className='h-10 w-10 relative'>
                                <Image src={"/painting.png"} alt={connector.name} fill={true} className='obeject-cover rounded-md' />
                            </div>
                            <div className='text-sm text-white text-center flex-1'>
                                {connector.name}
                            </div>
                        </div>
                    </li>)}
                </ul>
                <button className='w-full h-12 rounded-xl font-semibold bg-secondary'>Show More</button>
            </section>
            <p className='text-xs text-white text-center max-w-xs'>By using PowerDfi, you agree to our  Terms of Service and our Privacy Policy.</p>
        </div>
    </>

}

export default Connectors