import Link from 'next/link'
import React from 'react'
import { IoSettings } from "react-icons/io5";
import { IoPeopleOutline } from "react-icons/io5";



const HelpCenter = () => {
    return <section>
        <div className='flex justify-between'>
            <div className='w-1/3'>
                <h3 className='text-2xl text-white'>We Answer Your Every Need and Concern</h3>
            </div>
            <button className='px-8 ring-1 ring-primary/40 justify-center rounded-3xl bg-black-shade text-secondary text-sm h-10'>Get Help</button>
        </div>
        <div className='pt-2 pb-4 w-1/3'>
            <p className='text-xs text-white font-medium'>Your satisfaction and peace of mind are our priorities. We are committed to understanding and resolving your concerns</p>
        </div>

        <div className='grid grid-cols-2 gap-8 mt-8 w-1/2'>
            {
                [{
                    Icon: IoSettings,
                    title: "Account, Login & Billing",
                    note: "Troubleshooting guidance related to account login issues and billing concerns",

                }, {
                    Icon: IoPeopleOutline,
                    title: "Privacy and Policies",
                    note: "We address common queries and concerns regarding our privacy policies",
                }].map(({ Icon, title, note }) => <div key={title} className="aspect-square ring-1 rounded-lg ring-border p-4 flex flex-col gap-4 max-w-64">
                    <Icon className='text-secondary !text-4xl' />
                    <div className='flex flex-1 flex-col gap-2'>
                        <h3 className="text-white text-ms">{title}</h3>
                        <p className="text-white/80 text-xs tracking-tight">{note}</p>
                    </div>
                    <button className='w-full  flex items-center justify-center rounded-lg bg-border text-secondary text-sm h-10'>View</button>
                </div>)
            }


        </div>

    </section>
}

export default HelpCenter