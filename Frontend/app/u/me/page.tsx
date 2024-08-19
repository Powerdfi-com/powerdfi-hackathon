"use client"
import Image from 'next/image'
import React, { useContext, useEffect, useState } from 'react'
import { UserContext } from '@/app/context/context'
import UserAPI from '@/app/utils/apis/user'
import { Link } from '@/app/components/link'
import { Bar, BarChart, CartesianGrid, Legend, Line, LineChart, ResponsiveContainer, Tooltip, XAxis, YAxis } from 'recharts'
import Loading from '@/app/components/loading'
import { TStatAssets } from '@/app/utils/types'
import { MdArrowLeft, MdArrowRight } from 'react-icons/md'

const Dashboard = () => {
    const { user } = useContext(UserContext);
    const chart = [
        { "month": "Month 1", "sales": 752 },
        { "month": "Month 2", "sales": 434 },
        { "month": "Month 3", "sales": 899 },
        { "month": "Month 4", "sales": 123 },
        { "month": "Month 5", "sales": 678 },
        { "month": "Month 6", "sales": 302 },
        { "month": "Month 7", "sales": 912 },
        { "month": "Month 8", "sales": 456 },
        { "month": "Month 9", "sales": 789 },
        { "month": "Month 10", "sales": 345 },
        { "month": "Month 11", "sales": 645 },
        { "month": "Month 12", "sales": 564 }
    ]

    const { data, isPending } = UserAPI.getWalletDetails();
    const { data: activities, isPending: isPendingActivities } = UserAPI.getActivities();

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
    const { mutateAsync, data: topAssetsQuery } = UserAPI.getTopAssets();
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

    if (isPending || isPendingActivities) {
        return <Loading />
    }
    return <section className="flex flex-col gap-8 ml-12 pb-8">
        <div className='flex flex-col gap-6'>
            <h3 className='text-[30px] text-white-text leading-relaxed'>Welcome to your Token Dashboard</h3>
            <div className='flex gap-4'>
                <div className='px-6 h-80 flex-[4] flex items-center justify-center bg-black-shade rounded-2xl'>
                    <div className="flex flex-col py-6 px-4 gap-6 w-full">
                        <h3 className='w-2/3 text-[32px] text-white-text font-bold leading-relaxed'>Explore, Collect & Create your own Asset Token</h3>
                        <h3 className='text-base text-white-text leading-relaxed'>Explore base on asset type and trend</h3>

                        <div className='flex gap-6'>
                            <Link href="/u/h/explore" className="py-3 px-4 w-[201px] h-[43px] rounded-md bg-secondary text-sm font-semibold text-center">Explore</Link>
                            <Link href={"/u/me/create"} className="py-3 px-4 w-[201px] text-center h-[43px] rounded-md ring-1 ring-primary/40 text-sm text-secondary font-semibold">Create Asset</Link>
                        </div>
                    </div>
                </div>
                <div className="rounded-2xl bg-shade py-8 px-8 flex-[3] flex flex-col h-80">
                    <div className='flex items-center gap-4 py-8'>
                        <div className='h-14 w-14 rounded-full relative'>
                            <Image src={user.avatar} alt={user.username} fill className="rounded-full relative object-cover" />
                        </div>
                        <div className='flex flex-col justify-center'>
                            <p className="text-white/40 text-xs">Your total balance</p>
                            <p className="text-white-text text-2xl">{data?.data.balance}<span className='text-white/60 px-2'>PDfi</span></p>
                        </div>
                    </div>
                    <div className='flex justify-center items-center'>
                        <button className='w-full flex items-center justify-center ring-1 ring-primary/40 rounded-lg bg-black-shade text-secondary text-base font-bold py-3'>Get more PDfi</button>
                    </div>
                </div>
            </div>

            <div className='flex gap-6'>
                <div className='px-6 flex-[8] h-80 flex items-center justify-center bg-black-shade rounded-2xl'>
                    <ResponsiveContainer width="100%" height={300}>
                        <LineChart
                            data={chart}
                            margin={{
                                top: 20, right: 30, left: 20, bottom: 5,
                            }}
                        >
                            <XAxis dataKey="month" />
                            <YAxis />
                            <Tooltip />
                            <Line type="monotone" dataKey="sales" fill="#8884d8" />
                        </LineChart>
                    </ResponsiveContainer>
                </div>
                <div className="rounded-2xl flex-[4] bg-shade py-8 px-8 flex flex-col h-80 text-xl">
                    <div className='flex justify-center items-center'>
                        <p className='text-secondary text-xl'>Transaction History</p>
                    </div>
                    <table className='mt-4 table-fixed w-full'>
                        <thead>
                            <tr>
                                <th className='text-white text-md'>Action</th>
                                <th className='text-white text-md'>From/To</th>
                            </tr>
                        </thead>
                        {
                            activities?.data.activities.map((activity) => <tr key={activity.id}>
                                <td>
                                    <div className='flex flex-col gap-1'>
                                        <div className='flex gap-1'>
                                            <div className='text-white font-semibold capitalize'>{activity.action}</div>
                                        </div>
                                    </div>
                                </td>
                                <td>
                                    <div className='flex flex-col gap-1'>
                                        <div><span className='text-white/60'>From</span><span className='text-secondary'>{activity.fromUserId}</span></div>
                                        <div><span className='text-white/60'>To</span><span className='text-secondary'>{activity.toUserId}</span></div>
                                    </div>
                                </td>
                            </tr>)
                        }
                    </table>

                </div>
            </div>

            <div className='mt-8 w-[800px] h-[350px] ring-border ring-1 justify-start rounded-lg p-6 flex flex-col'>
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
        </div>

    </section>
}

export default Dashboard


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