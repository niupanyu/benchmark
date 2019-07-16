#include <pthread.h>
#include <iostream>
#include <vector>

void start_routine(void *arg)
{
    int i;
    i += i;
    return;
}

int main()
{
    std::vector< pthread_t > vTid;
    std::size_t count = 10000;
    pthread_t  tid;
    for(std::size_t i = 0 ; i < count; i++)
    {
        int ret = pthread_create(&tid, 0, (void*(*)(void*))&start_routine, (void *)&i);
        vTid.push_back(tid);
    }

    for(std::size_t i = 0; i < vTid.size(); i++)
    {
        pthread_join(vTid[i],NULL);
    }
    return 0;
}
