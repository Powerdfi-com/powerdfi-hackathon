"use client"
import Loading from '@/app/components/loading';
import AdminAPI from '@/app/utils/apis/admin'
import { months } from '@/app/utils/func';
import React from 'react'
import { CartesianGrid, Legend, Line, LineChart, ResponsiveContainer, Tooltip, XAxis, YAxis } from 'recharts';

const SiteAnalytics = () => {
    const { data: creatorsSurvey } = AdminAPI.getCreatorsSurvey();
    const { data: categoriesSurvey } = AdminAPI.getCategoriesSurvey();
    const { data: assetsSurvey } = AdminAPI.getAssetsSurvey();
    const colors = [
        "#1bdefd",
        "#c02cf4",
        "#303030",
        "#1BDEFD80",
        "#f60",
        "white",
        "blue",
        "#f80"
    ]
    const tailwindColors = [
        "bg-[#1bdefd]",
        "bg-[#c02cf4]",
        "bg-[#303030]",
        "bg-[#1BDEFD80]",
        "bg-[#f60]",
        "bg-[white]",
        "bg-[blue]",
        "bg-[#f80]"
    ]
    if (!creatorsSurvey || !categoriesSurvey || !assetsSurvey) return <Loading />
    return <section>
        <h3 className='text-white text-2xl'>Site Analytics</h3>
        <div className='w-[800px] rounded-xl ring-1 mt-6 ring-border flex flex-col p-4 gap-6'>
            <h4 className='text-white text-lg'>Creators Survery</h4>
            <div className='flex-1'>
                <ResponsiveContainer width="100%" height={400}>
                    <LineChart
                        data={creatorsSurvey.data.map((_data) => {
                            return { ..._data, monthName: months[_data.month - 1] }
                        })}
                        margin={{
                            top: 5, right: 30, left: 20, bottom: 5,
                        }}
                    >
                        <CartesianGrid strokeWidth={0.2} />
                        <XAxis dataKey="monthName" />
                        <YAxis />
                        <Tooltip />
                        <Line type="linear" dataKey="heritageUsersCount" stroke="#1bdefd" />
                        <Line type="linear" dataKey="newUsersCount" stroke="#c02cf4" />
                    </LineChart>
                </ResponsiveContainer>
            </div>
            <div className='flex justify-center w-full gap-8'>
                <div className='flex items-center gap-4'>
                    <div className='bg-secondary rounded-md w-[40px] h-[30px]'></div>
                    <span className='text-sm text-white'>Heirtage Creators</span>
                </div>
                <div className='flex items-center gap-4'>
                    <div className='bg-primary rounded-md w-[40px] h-[30px]'></div>
                    <span className='text-sm text-white'>New Creators</span>
                </div>
            </div>
        </div>
        <div className='w-[800px] rounded-xl ring-1 mt-6 ring-border flex flex-col p-4 gap-6'>
            <h4 className='text-white text-lg'>Creators survey by asset verification</h4>
            <div className='flex-1'>
                <ResponsiveContainer width="100%" height={400}>
                    <LineChart
                        data={assetsSurvey.data.map((_data) => {
                            return { ..._data, monthName: months[_data.month - 1] }
                        })}
                        margin={{
                            top: 5, right: 30, left: 20, bottom: 5,
                        }}
                    >
                        <CartesianGrid strokeWidth={0.2} />
                        <XAxis dataKey="monthName" />
                        <YAxis />
                        <Tooltip />
                        <Line type="linear" dataKey="verifiedAssetsCount" stroke="#1bdefd" />
                        <Line type="linear" dataKey="unverifiedAssetsCount" stroke="#c02cf4" />
                    </LineChart>
                </ResponsiveContainer>
            </div>
            <div className='flex justify-center w-full gap-8'>
                <div className='flex items-center gap-4'>
                    <div className='bg-secondary rounded-md w-[40px] h-[30px]'></div>
                    <span className='text-sm text-white'>Verified Assets</span>
                </div>
                <div className='flex items-center gap-4'>
                    <div className='bg-primary rounded-md w-[40px] h-[30px]'></div>
                    <span className='text-sm text-white'>Rejected Assets</span>
                </div>
            </div>
        </div>
        <div className='w-[800px] rounded-xl ring-1 mt-6 ring-border flex flex-col p-4 gap-6'>
            <h4 className='text-white text-lg'>Site Statistics by asset type</h4>
            <div className='flex-1'>
                <ResponsiveContainer width="100%" height={300}>
                    <LineChart
                        data={categoriesSurvey.data.map((_data) => {
                            return { ..._data, monthName: months[_data.month - 1] }
                        })}
                        margin={{
                            top: 5, right: 30, left: 20, bottom: 5,
                        }}
                    >
                        <CartesianGrid strokeWidth={0.2} />
                        <XAxis dataKey="monthName" />
                        <YAxis />
                        <Tooltip />
                        {
                            Object.keys(categoriesSurvey.data[0]).filter((category) => category !== "month").map((category, index) => <Line key={category} type="linear" dataKey={category} stroke={colors[index]} />)
                        }

                    </LineChart>
                </ResponsiveContainer>
            </div>
            <div className='grid max-w-xs grid-cols-2 justify-center w-full gap-2 self-center'>
                {
                    Object.keys(categoriesSurvey.data[0]).filter((category) => category !== "month").map((category, index) => <div key={category} className='flex items-center gap-4'>
                        <div className={`${tailwindColors[index]} rounded-md w-[40px] h-[30px]`}></div>
                        <span className='text-sm text-white upp'>{category}</span>
                    </div>)
                }

            </div>
        </div>

    </section>
}

export default SiteAnalytics