"use client"
import Asset from '@/app/components/asset'
import Error from '@/app/components/error'
import Loading from '@/app/components/loading'
import { UserContext } from '@/app/context/context'
import UserAPI from '@/app/utils/apis/user'
import React, { useContext } from 'react'

const Portfolio = () => {
    const { user } = useContext(UserContext)
    const { data, isPending } = UserAPI.getPortfolio();
    return <section>
        <div className='flex justify-between'>
            <h3 className='text-2xl text-white leading-relaxed'>Your Asset Portfolio</h3>
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
                {data.data.map((asset) => <Asset key={asset.id} isListed={false} asset={asset} />)}
            </ul> : <Error />)
        }
    </section>
}

export default Portfolio