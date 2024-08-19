"use client"
import Asset from '@/app/components/asset'
import Error from '@/app/components/error'
import { Link } from '@/app/components/link'
import Loading from '@/app/components/loading'
import { UserContext } from '@/app/context/context'
import UserAPI from '@/app/utils/apis/user'
import React, { useContext } from 'react'
import { GoArrowSwitch } from 'react-icons/go'

const Lisitng = () => {
    const { user } = useContext(UserContext);
    const { data, isPending } = UserAPI.getListedAssets(user.id);
    return <section>
        <div className='flex justify-between'>
            <h3 className='text-2xl text-white leading-relaxed'>Your Listed Tokens</h3>
            <Link href="/u/me/portfolio" className='rounded-full flex items-center bg-black-shade'><div className='h-8 w-8 rounded-full bg-secondary flex items-center justify-center'><GoArrowSwitch /></div><span className='text-white text-sm px-3'>Create New Token Lisiting</span></Link>
        </div>
        <div className='mt-6 border-y border-border py-5 flex items-center justify-between'>
            <div className='flex items-center gap-1'>
                <span className='text-xs text-white'>Date Range</span>
                <input type='date' className='text-xs text-white/60 bg-shade p-1.5 rounded-md' />
                <input type='date' className='text-xs text-white/60 bg-shade p-1.5 rounded-md' />
            </div>
            <input className='bg-transparent rounded-lg text-white text-sm p-1.5 border-border border w-full max-w-sm' type='search' />
            <select className='bg-shade p-1.5 rounded-md text-white text-xs text'>
                <option>Sort By</option>
            </select>
        </div>
        {
            isPending ? <Loading /> : (data ? <ul className='grid grid-cols-3 gap-6 mt-8'>
                {data.data.map((asset) => <Asset key={asset.id} asset={asset} isListed={true} />)}
            </ul> : <Error />)
        }
    </section>
}

export default Lisitng