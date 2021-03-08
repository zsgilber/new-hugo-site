# for my programming resume, only choose work positions I started after QDI (Jan 31, 2016)
.work |= map(select(.startDate >= "2016-01-31"))