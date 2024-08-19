"use client"
import Error from '@/app/components/error';
import { Link } from '@/app/components/link';
import Loading from '@/app/components/loading';
import AssetAPI from '@/app/utils/apis/asset'
import Image from 'next/image';
import React from 'react'

const Chains = () => {
    const { data, isPending } = AssetAPI.getChains();
    return <section className='flex flex-col gap-4'>
        <h3 className='text-2xl text-white leading-relaxed'>Choose Blockchain</h3>
        <p className='text-white text-sm'>Choose the most suitable blockchain for your needs.</p>
        {isPending ? <Loading /> : (data ? <ul className='bg-blue-shade rounded-xl grid grid-cols-3 p-8 gap-16'>
            {
                data.data.map((chain) => <li key={chain.id}>
                    <Link href={`/u/me/create/chains/${chain.id}`}><div className='w-full max-w-sm aspect-square flex flex-col items-center justify-center gap-4 ring-1 ring-white/5 rounded-xl'>
                        <div className='h-16 w-16 relative rounded-xl'>
                            <Image src={chain.logo} alt={chain.name} fill={true} className='rounded-xl object-cover' />
                        </div>
                        <div className='text-white text-sm'>{chain.name}</div>
                    </div></Link>
                </li>)
            }
        </ul> : <Error />)}
    </section>
}

export default Chains