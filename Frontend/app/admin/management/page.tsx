"use client"
import Error from '@/app/components/error';
import Loading from '@/app/components/loading';
import AdminAPI from '@/app/utils/apis/admin';
import { useRouter } from 'next/navigation';
import React from 'react'

const Management = () => {
    const { data, isPending } = AdminAPI.getAssets();
    const router = useRouter()
    return <section>
        <div className='flex justify-between'>
            <h3 className='text-2xl text-white leading-relaxed'>Asset Management</h3>
        </div>
        <div className='border-y border-border h-14 mt-4 flex items-center justify-between gap-8'>
            <div className="flex gap-2 items-center rounded-md bg-black-shade px-2 py-0.5">
                <span className='text-white text-sm'>Date Range</span>
                <input type='date' className="bg-transparent  rounded-md text-white/40 text-xs" />
                <input type='date' className="bg-transparent  rounded-md text-white/40 text-xs" />
            </div>
            <input type="search" placeholder='Search' className="bg-transparent px-2 py-0.5 flex-1 rounded-md ring-1 ring-border" />
            <select className="bg-black-shade text-white px-2 py-0.5 rounded-md">
                <option>Sort By</option>
            </select>
        </div>
        <table className="w-full">
            <thead>
                <tr className='text-md border-b border-border text-white h-14'>
                    <th className='text-left'>Asset</th>
                    <th className='text-left'>Asset Type</th>
                    <th className='text-left'>Floor Price</th>
                    <th className='text-left'>Total Supply</th>
                    <th className='text-left'>Status</th>
                </tr>
            </thead>
            {
                isPending ? <Loading className="h-0" /> : (data ? <tbody>
                    {
                        data.data.map((asset) => <tr key={asset.id} className='text-md border-b border-border text-white h-14 cursor-pointer' onClick={() => router.push(`/admin/management/${asset.id}`)}>
                            <td className='text-left'>{asset.name}</td>
                            <td className='text-left'>{asset.symbol}</td>
                            <td className='text-left'>{asset.floorPrice}</td>
                            <td className='text-left'>{asset.totalSupply}</td>
                            <td className='text-left'>
                                {
                                    asset.status === "verified" && <button className='h-8 px-8 bg-green-950 text-green-100 text-sm ring-1 ring-green-600 rounded-md'>Verified</button>
                                }
                                {
                                    asset.status === "rejected" && <button className='h-8 px-8 bg-red-950 text-red-100 text-sm ring-1 ring-red-600 rounded-md'>Rejected</button>
                                }
                                {
                                    asset.status === "unverified" && <button className='h-8 px-8 bg-blue-950 text-blue-100 text-sm ring-1 ring-blue-600 rounded-md'>Awating Verification</button>
                                }
                            </td>
                        </tr>)
                    }
                </tbody> : <Error />)
            }
        </table>
    </section>
}

export default Management