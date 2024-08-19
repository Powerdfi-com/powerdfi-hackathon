"use client"
import Loading from '@/app/components/loading';
import { GhostAPI } from '@/app/utils/apis/ghost';
import { dateFromISO, formatDate } from '@/app/utils/func';
import Image from 'next/image';
import React from 'react'
import { IoPersonSharp } from 'react-icons/io5';

const GhostBlog = ({ slug }: { slug: string }) => {
    const { data, isLoading } = GhostAPI.fetchBlogBySlug(slug);
    if (isLoading) {
        return <Loading />
    }
    return <main className="bg-gradient-radial from-black to-primary/5 px-8 sm:px-12 md:px-24 py-12">
        <h3 className="text-white text-[36px] text-center font-semibold"><q>{data?.title}</q></h3>
        <div className="flex justify-center items-center text-[14px] text-text-ash space-x-8 py-4">
            {
                data?.authors && <div className="flex justify-center items-center space-x-3">
                    <div className="rounded-full bg-white/40 p-2">
                        <IoPersonSharp className="!text-2xl text-black" />
                    </div>
                    <p className="text-white/80">{data!.authors[0].name}</p>
                </div>
            }
            <p className="text-white/80">{formatDate(dateFromISO(data!.updated_at!))}</p>
        </div>
        <div className="w-full h-[560px] relative mt-16 rounded-2xl">
            <Image src={data?.feature_image!} alt={data?.slug!} fill className="object-cover rounded-2xl" />
        </div>
        <section className="text-white mt-12" dangerouslySetInnerHTML={{ __html: data!.html! }}></section>
    </main>
}

export default GhostBlog