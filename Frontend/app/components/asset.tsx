import React from 'react'
import { TAsset } from '../utils/types'
import { Link } from '@/app/components/link'
import { BsAward } from "react-icons/bs";
import Image from 'next/image';
import Verified from './verified';



const Asset = ({ asset, isListed }: { asset: TAsset, isListed: boolean }) => {
    return <div key={asset.id} className="aspect-square ring-1 rounded-lg ring-border p-4 flex flex-col gap-4 max-w-64">
        <div className='flex flex-1 flex-col gap-2'>
            <div className='flex gap-4 '>
                <div className='flex-1 flex flex-col gap-2'>
                    <h4 className="text-secondary text-xs">{asset.symbol}</h4>
                    <p className="text-white text-3xl">{asset.name}</p>
                </div>
                {
                    asset.status === "verified" && <Verified />
                }
            </div>
            <div className='rounded-full mt-4 ring-1 px-4 py-0.5 ring-primary/40 text-xs text-secondary w-fit'>{asset.category}</div>
        </div>
        {
            isListed ? <Link href={`/u/me/portfolio/${asset.id}`} className='w-full  flex items-center justify-center rounded-lg bg-black-shade text-secondary text-sm h-10'>Manage Token</Link> : <Link href={`/u/me/listing/create/${asset.id}/${asset.blockchain}`} className='w-full  flex items-center justify-center rounded-lg bg-black-shade text-secondary text-sm h-10'>List Asset</Link>
        }
    </div>
}

export default Asset