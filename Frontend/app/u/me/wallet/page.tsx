"use client"
import Link from 'next/link'
import React, { useContext } from 'react'
import { useWalletInfo } from "@web3modal/wagmi/react";
import { UserContext } from '@/app/context/context'
import { IoCopy } from 'react-icons/io5'
import Image from 'next/image'
import UserAPI from '@/app/utils/apis/user';

const Wallet = () => {
    const { user, address, copyAddress } = useContext(UserContext);
    const { data } = UserAPI.getWalletDetails();
    return <section className="pb-12">
        <div className='flex justify-between'>
            <h3 className='text-2xl text-white leading-relaxed'>Wallet</h3>
        </div>
        <div className='ring-1 ring-primary/20 px-6 flex justify-center items-center  gradient rounded-lg w-[748px] ml-10 mt-12 h-[508px]'>
            <div className="flex flex-col items-center gap-6 w-[542px]">
                <div className="flex gap-6 items-center w-full">
                    <div className="h-16 w-16 flex items-center justify-center ring-1 ring-border rounded-md">
                        <div className="relative h-10 w-10">
                            <Image src={"/icon.png"} alt="user" className="object-cover" fill={true} />
                        </div>
                    </div>
                    <div className="flex flex-col flex-1">
                        <div className="text-xs text-white/40 font-light">Your wallet address</div>
                        <div className="text-white text-2xl max-w-xs w-full overflow-hidden text-ellipsis">{address}</div>
                    </div>
                    <button className="text-white !font-semibold !text-lg" onClick={copyAddress}><IoCopy /></button>
                </div>
                {
                    data && <div className='text-4xl ring-1 flex items-center justify-center h-[96px] ring-border rounded-md w-full'>
                        <span className='text-white'>{data.data.balance} <span className='text-secondary'>USD</span></span>
                    </div>
                }
                <button className="py-2 w-full rounded-md bg-secondary text-sm font-semibold max-w-[289px]">Add Funds</button>
            </div>
        </div>


    </section>
}

export default Wallet