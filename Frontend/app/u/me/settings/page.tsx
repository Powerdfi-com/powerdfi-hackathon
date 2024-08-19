import { Link } from '@/app/components/link'
import React from 'react'
import { GoArrowSwitch } from 'react-icons/go'

const Settings = () => {
    return <section>
        <div className='flex justify-between'>
            <h3 className='text-2xl text-white leading-relaxed'>Your Listed Tokens</h3>
            <Link href="/u/me/listing/new" className='rounded-full flex items-center bg-black-shade'><div className='h-8 w-8 rounded-full bg-secondary flex items-center justify-center'><GoArrowSwitch /></div><span className='text-white text-sm px-3'>Create New Token Lisiting</span></Link>
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
        <ul className='grid grid-cols-3 gap-6 mt-8'>
            {[1, 2, 3, 4, 5].map((e) => <div key={e} className="aspect-square ring-1 rounded-lg ring-border p-4 flex flex-col gap-4 max-w-64">
                <div className='flex flex-1 flex-col gap-2'>
                    <h4 className="text-secondary text-xs">Elton Ave</h4>
                    <p className="text-white">Elton Ave</p>
                    <div className='rounded-full mt-4 ring-1 px-4 py-0.5 ring-primary/40 text-xs text-secondary w-fit'>Real Estate</div>
                </div>
                <Link href={"/u/me/listing/index"} className='w-full  flex items-center justify-center rounded-lg bg-black-shade text-secondary text-sm h-10'>Preview Token</Link>
            </div>)}
        </ul>
    </section>
}

export default Settings