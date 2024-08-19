"use client"
import Error from '@/app/components/error';
import Loading from '@/app/components/loading';
import { AssetContext } from '@/app/context/context';
import AssetAPI from '@/app/utils/apis/asset';
import Image from 'next/image'
import React, { useContext } from 'react'

const History = ({ params }: { params: { assetId: string } }) => {
    const asset = useContext(AssetContext);
    const { data, isPending } = AssetAPI.getActivities({
        assetId: params.assetId
    })
    return <section>
        <div className='flex justify-between h-14 border-y border-border items-center'>
            <div className='bg-black-shade p-1 rounded-md text-white text-xs'>Latest Transaction</div>
            <select className='bg-black-shade p-1 rounded-md text-white text-xs'>
                <option>Sort By</option>
            </select>
        </div>
        <table className='w-full'>
            <thead>
                <tr className='h-14 text-[20px] text-white'>
                    <th className='text-left'>Action</th>
                    <th className='text-left'>Assets</th>
                    <th className='text-left'>Quantity</th>
                    <th className='text-left'>Price</th>
                    <th className='text-left'>From/To</th>
                </tr>
            </thead>
            {
                isPending ? <Loading /> : (data ? data.data.map((transaction) => <tr key={transaction.id} className='text-[20px] text-white/80 h-14 text-left'>
                    <td>{transaction.action}</td>
                    <td>{transaction.assetName}</td>
                    <td>{transaction.quantity}</td>
                    <td>{transaction.price}</td>
                    <td><div className='flex flex-col gap-1'>
                        <div>From: {transaction.fromUserId}</div>
                        <div>To: {transaction.toUserId}</div>
                    </div></td>
                </tr>) : <Error />)
            }
        </table>
    </section>
}

export default History