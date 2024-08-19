"use client"
import React, { useEffect, useState } from 'react'
import Image from 'next/image';
import AdminAPI from '@/app/utils/apis/admin';
import StatsAPI from '@/app/utils/apis/stats';
import { TStatAssets } from '@/app/utils/types';
import { MdArrowLeft, MdArrowRight } from 'react-icons/md';
import { Line, LineChart } from 'recharts';
import Loading from '@/app/components/loading';

const UsersAnalytics = () => {
    const { data } = AdminAPI.getStats();
    const hours = [
        { text: "1hr", value: 1 },
        { text: "24hrs", value: 24 },
        { text: "7days", value: 7 },
        { text: "30days", value: 30 },
    ]
    const [totalStats, setTotalStats] = useState(0)
    const [hour, setHour] = useState(1);
    const [page, setpage] = useState(1);
    const [stats, setStats] = useState<TStatAssets["assets"]>([])
    const { isPending, mutateAsync, data: topAssetsQuery } = StatsAPI.getTopAssets();
    useEffect(() => {
        const fetchData = async () => {
            await mutateAsync({ page, size: 5, range: hour }).then((res) => {
                setStats((_stats) => [..._stats, ...res.data.assets.filter((stat) => !_stats.find((_stat) => _stat.assetId === stat.assetId))]);
                setTotalStats(res.data.total);
            })
        }
        fetchData()
    }, [page, mutateAsync, hour])
    const totalPages = Math.ceil(totalStats / 5);

    return <section>
        <div className='flex justify-between'>
            <h3 className='text-2xl text-white leading-relaxed'>Users Analytics</h3>
        </div>

        <div className='mt-8 w-[800px] h-[350px] ring-border ring-1 rounded-sm p-6'>
            <h3 className='text-white/60 text-[24px]'>Users Statistics</h3>
            <div className='flex-1 w-full flex'>
                <div className='flex flex-col gap-8 mt-8'>
                    <div className='flex flex-col'>
                        <div className='text-[24px] text-white'>
                            Creators
                        </div>
                        <div className='text-white/60 text-[16px]'>
                            Last week raised by <span className='text-secondary'>{data?.data.percentageChangeCreators}</span>
                        </div>
                    </div>
                    <div className='flex flex-col'>
                        <div className='text-[24px] text-white'>
                            Users
                        </div>
                        <div className='text-white/60 text-[16px]'>
                            Last week raised by <span className='text-secondary'>{data?.data.percentageChangeUsers}</span>
                        </div>
                    </div>
                </div>
                <div className='flex-1 h-[250px] relative '>
                    <Image src="/map.png" alt="map" fill={true} />
                </div>
            </div>
        </div>

        <div className='mt-8 w-[800px] h-[350px] ring-border ring-1 justify-start rounded-sm p-6 flex flex-col'>
            <div className='w-full flex justify-between'>
                <div className='text-white/60 text-[24px]'>Top Asset Sales</div>
                <ul className='flex gap-2 rounded-lg bg-black-shade items-center'>
                    {hours.map(({ text, value }) => <li key={value}><div className={(value === hour ? 'bg-blue-shade rounded-lg' : '') + (' text-white py-2 px-4 text-xs cursor-pointer')} onClick={() => {
                        setHour(value);
                        setpage(1)
                    }}>{text}</div></li>)}
                </ul>
            </div>
            <div className='w-full flex-1 overflow-auto'>
                <table className='table-fixed w-full mt-6'>
                    <thead>
                        <tr>
                            <th className='text-white/60 text-[20px] text-left'>Asset</th>
                            <th className='text-white/60 text-[20px] text-left'>Floor Price</th>
                            <th className='text-white/60 text-[20px] text-left'>Last 24h</th>
                            <th className='text-white/60 text-[20px] text-left'>Volume</th>
                        </tr>
                    </thead>
                    {topAssetsQuery ? topAssetsQuery?.data.assets.map((asset) => <Asset asset={asset} key={asset.assetId} />) : <Loading className="!w-[800px] h-full" />}
                </table>
            </div>
            <div className='flex justify-end'>
                <div className='flex gap-2'>
                    <button disabled={page === 1} onClick={() => setpage((_page) => _page - 1)}>
                        <MdArrowLeft className='text-secondary !text-xl' />
                    </button>
                    <span className='text-white/60 text-sm'>{page > totalPages ? totalPages : page}/{totalPages}</span>
                    <button disabled={page >= totalPages} onClick={() => setpage((_page) => _page + 1)}><MdArrowRight className='text-secondary !text-xl' /></button>
                </div>
            </div>
        </div>
    </section>
}

export default UsersAnalytics



const Asset = ({ asset }: { asset: TStatAssets["assets"][0] }) => {
    return <tr className='text-white text-[20px] h-[68px]'>
        <td>{asset.name}</td>
        <td>{asset.floorPrice}</td>
        <td className='flex items-center gap-2'><Chart asset={asset} />{asset.percentageChange}</td>
        <td>{asset.volume}</td>
    </tr>
}

const Chart = ({ asset }: { asset: TStatAssets["assets"][0] }) => {
    const data = asset.priceChanges?.length === 0 ? [
        { name: 'Point 1', value: 5 },
        { name: 'Point 2', value: 5 },
    ] : [
        { name: 'Point 1', value: 5 },
        { name: 'Point 2', value: 5 },
    ];
    return <LineChart width={100} height={68} data={data}>
        <Line dot={false} type="linear" dataKey="value" stroke="#1bdefd" yAxisId={0} />
    </LineChart>
}