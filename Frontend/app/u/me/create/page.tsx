"use client"
import Asset from '@/app/components/asset'
import { Link } from '@/app/components/link'
import { UserContext } from '@/app/context/context'
import UserAPI from '@/app/utils/apis/user'
import React, { useContext } from 'react'
import { FaPlus } from 'react-icons/fa6'

const Create = () => {
    const { user } = useContext(UserContext)
    const { isPending, data } = UserAPI.getCreatedAssets(user.id);
    return <section className="">
        <h3 className='text-xl text-white leading-relaxed'>Your asset tokens</h3>
        <ul className='grid grid-cols-3 gap-6 mt-8'>
            <Link href={"/u/me/create/chains"} className="aspect-square ring-1 rounded-lg ring-border p-4 max-w-64 flex flex-col gap-4 justify-center items-center">
                <FaPlus className='text-white/60 !text-2xl' />
                <span className="text-white text-sm">New Asset Token</span>
            </Link >
            {data && data.data.map((asset) => <Asset key={asset.id} asset={asset} isListed={false} />)}
        </ul>
    </section>
}

export default Create